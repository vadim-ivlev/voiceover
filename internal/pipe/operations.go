package pipe

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/audio"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/pkg/epubs"
	"github.com/vadim-ivlev/voiceover/pkg/sound"
	"github.com/vadim-ivlev/voiceover/pkg/texts"
	"github.com/vadim-ivlev/voiceover/pkg/translator"
	"golang.org/x/exp/rand"
)

// Nap - Sleep for a random duration between 0 and NapTime milliseconds
func Nap() {
	if config.Params.NapTime == 0 {
		return
	}
	// random duration between 0 and NapTime milliseconds
	napDuration := time.Duration(rand.Intn(config.Params.NapTime)) * time.Millisecond
	time.Sleep(napDuration)
}

// LogJob - logs the job information
func LogJob(job Job, msg string) {
	log.Info().Msgf("Job %d %s", job.ID, msg)
	// log.Info().Msg(job.String())
}

// getTextJobOperation0 - returns the textJob operation
func getTextOperation0(textLines []string) JobFunction {
	// adds text to the job, and assigns a voice
	return func(job Job) (Job, error) {
		if len(textLines) < job.ID+1 {
			return job, fmt.Errorf("Job %d: No text for the job", job.ID)
		}
		job.Results.Text = strings.TrimSpace(textLines[job.ID])
		return job, nil
	}
}

// getTextJobOperation - returns the textJob operation
func getTextOperation(textLines []epubs.EpubTextLine) JobFunction {
	// adds text to the job, and assigns a voice
	return func(job Job) (Job, error) {
		if len(textLines) < job.ID+1 {
			return job, fmt.Errorf("Job %d: No text for the job", job.ID)
		}
		job.Results.Epub = textLines[job.ID]
		job.Results.Text = strings.TrimSpace(textLines[job.ID].Text)
		job.Results.Html = textLines[job.ID].Html
		return job, nil
	}
}

// translateTextOperation - translates the text
func translateTextOperation(job Job) (Job, error) {
	// check if the text is empty or we don't need to translate
	if job.Results.Text == "" || config.Params.TranslateTo == "" {
		job.Results.TranslatedText = strings.TrimSpace(job.Results.Text)
		return job, nil
	}

	// translate the text
	translatedText, err := translator.TranslateText(config.Params.OpenaiAPIURL, config.Params.ApiKey, config.Params.TranslateTo, translator.TextTranslInstructions, job.Results.Text)
	if err != nil {
		return job, err
	}
	// save the translated text to the job
	job.Results.TranslatedText = translatedText

	// translate the html
	TranslatedHtml, err := translator.TranslateText(config.Params.OpenaiAPIURL, config.Params.ApiKey, config.Params.TranslateTo, translator.HtmlTranslInstructions, job.Results.Html)
	if err != nil {
		return job, err
	}
	// save the translated text to the job
	job.Results.TranslatedHtml = TranslatedHtml

	return job, nil
}

// soundOperation - generates a sound file and a text file for the job
func soundOperation(job Job) (Job, error) {
	job.Results.TextFile = fmt.Sprintf("%s/%08d.txt", config.Params.TextsDir, job.ID)
	job.Results.AudioFile = fmt.Sprintf("%s/%08d.mp3", config.Params.SoundsDir, job.ID)

	// select voice
	if len(job.Results.TranslatedText) > 0 {
		job.Results.Voice = audio.NextVoice()
	}

	var err error

	// if voice is empty, generate silence
	if job.Results.Voice == "" {
		err = sound.GenerateSilenceMP3(config.Params.Pause, job.Results.AudioFile)
	} else {
		err = audio.GenerateSpeechMP3(config.Params.Speed, job.Results.Voice, job.Results.TranslatedText, job.Results.AudioFile)
	}
	if err != nil {
		return job, err
	}

	err = texts.SaveTextFile(job.Results.TextFile, job.Results.TranslatedText)
	if err != nil {
		return job, err
	}

	return job, nil
}

// toArray - moves the jobs from the channel to an array
// Parameters:
// jobsChan: the channel containing the jobs
// Returns:
// an array of jobs
func toArray(jobsChan chan Job) []Job {
	jobsArray := []Job{}
	for job := range jobsChan {
		jobsArray = append(jobsArray, job)
	}
	// restore the initial order
	sort.Slice(jobsArray, func(i, j int) bool {
		return jobsArray[i].ID < jobsArray[j].ID
	})
	return jobsArray
}

// DoPipeline - starts the pipeline processing.
// Parameters:
//   - task: the task to process
//
// Returns:
//   - an array of processed jobs
//   - the number of jobs
//   - the base name of the output files
//   - an error if any
func DoPipeline(task Task) (doneJobs []Job, numJobs int, outputBaseName string, err error) {

	// get slice of epub text lines and the base name for the output files
	epubTextLines, outputBaseName, err := getEpubTextLines()
	if err != nil {
		return
	}

	numJobs = len(epubTextLines)

	// create an array of jobs
	jobs := newJobsArray(numJobs)
	// load the previous jobs to the new jobs array
	loadPreviuosJobs(jobs, task)

	// Create channels to pass the jobs between the workers
	jobsChan := make(chan Job, numJobs)
	textChan := make(chan Job, numJobs)
	translChan := make(chan Job, numJobs)
	// soundChan must be buffered because its jobs will be sorted after it is closed
	soundChan := make(chan Job, numJobs)

	// Fill the jobs channel with the jobs from a newly created array
	go toChannel(jobs, jobsChan)

	// add a text to each job and assign a voice
	// doTeamWork(4, "Text", "T", getTextOperation(epubTextLines), jobsChan, textChan)
	go DoWork(nil, "Add Text", "T", getTextOperation(epubTextLines), jobsChan, textChan)

	// translate the text.
	// go DoWork(nil, "Translate", "Tr", translateTextOperation, textChan, translChan)
	doTeamWork(10, "Translate", "Tr", translateTextOperation, textChan, translChan)

	// generate sound file for each job. Fan-out.
	doTeamWork(10, "Sound", "S", soundOperation, translChan, soundChan)

	// gatther the jobs into an array. Fan-in.
	doneJobs = toArray(soundChan)
	return doneJobs, numJobs, outputBaseName, nil
}

// CreateFileList - creates a file list of audio files for concatenation with ffmpeg
// File list example:
// file 'file1.mp3'
// file 'file2.mp3'
// file 'file3.mp3'
// Parameters:
// jobs: the jobs to process
func CreateFileList(jobs []Job) (err error) {

	RemoveFileListFile()

	fileList := ""
	for _, job := range jobs {
		fileList += fmt.Sprintf("file '%s'\n", job.Results.AudioFile)
	}
	err = texts.SaveTextFile(config.Params.FileListFileName, fileList)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFileListFile() {
	err := os.RemoveAll(config.Params.FileListFileName)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
}

// CreateOutputMP3 - joins the processed jobs into one mp3 file.
func CreateOutputMP3(processedJobs []Job, outputBaseName string) (outMP3File string, err error) {
	// create file list of audio files for concatenation
	err = CreateFileList(processedJobs)
	if err != nil {
		return
	}

	// remove the output mp3 file if it exists
	err = os.Remove(outputBaseName + ".mp3")
	if err != nil {
		log.Info().Msgf("Failed to delete the output file: %v", err)
	}

	// concatenate the audio files into one
	outMP3File = outputBaseName + ".mp3"
	err = audio.ConcatenateMP3Files(config.Params.FileListFileName, outMP3File)
	return
}

// CreateOutputText - joins the processed jobs into one text file.
func CreateOutputText(processedJobs []Job, outputBaseName string) (outTextFile string, err error) {
	// remove the output text file if it exists
	err = os.Remove(outputBaseName + ".txt")
	if err != nil {
		log.Info().Msgf("Failed to delete the output file: %v", err)
	}

	// get the text lines from the processed jobs
	textLines := []string{}
	for _, job := range processedJobs {
		textLines = append(textLines, job.Results.TranslatedText)
	}

	// write a text file with processed lines
	outTextFile = outputBaseName + ".txt"
	err = texts.SaveTextFileLines(outTextFile, textLines)
	if err != nil {
		return
	}

	return
}

// CreateOutputEpub - joins the processed jobs into one epub file.
func CreateOutputEpub(inputEpubFile string, processedJobs []Job, outputBaseName string) (outEpubFile string, err error) {
	// remove the output text file if it exists
	err = os.Remove(outputBaseName + ".epub")
	if err != nil {
		log.Info().Msgf("Failed to delete the output file: %v", err)
	}

	// get the text lines from the processed jobs
	epubLines := []epubs.EpubTextLine{}
	for _, job := range processedJobs {
		line := job.Results.Epub
		line.Text = job.Results.TranslatedText
		line.Html = job.Results.TranslatedHtml
		epubLines = append(epubLines, line)
	}

	// write an epub file with processed lines
	outEpubFile = outputBaseName + ".epub"

	if inputEpubFile == outEpubFile {
		err = fmt.Errorf("input and output files are the same: %s", inputEpubFile)
		return
	}

	// err = epubs.SaveEpubTextLines(outEpubFile, epubLines)
	err = epubs.RepackEpub(inputEpubFile, outEpubFile, epubLines)
	if err != nil {
		return
	}

	return
}

// LoadJSONFile - loads a JSON file into a structure
// Parameters:
// fileName: the name of the file
// v: the structure to load the file into
func LoadJSONFile(fileName string, v any) (err error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return
	}
	return nil
}

// getEpubTextLinesFromTextFile reads lines from a text file and converts them into a slice of EpubTextLine structs.
//
// Parameters:
//   - fileName: The name of the text file to read from.
//   - startIndex: The starting index of the lines to read.
//   - endIndex: The ending index of the lines to read.
//
// Returns:
//   - epubTextLines: A slice of EpubTextLine structs containing the text lines and their metadata.
//   - start: The actual starting index of the lines read.
//   - end: The actual ending index of the lines read.
//   - err: An error object if an error occurred while reading the file, otherwise nil.
func getEpubTextLinesFromTextFile(fileName string, startIndex, endIndex int) (epubTextLines []epubs.EpubTextLine, start, end int, err error) {
	epubTextLines = []epubs.EpubTextLine{}
	lines, start, end, err := texts.GetTextFileLines(fileName, startIndex, endIndex)
	if err != nil {
		return
	}
	for i, line := range lines {
		epubTextLines = append(epubTextLines, epubs.EpubTextLine{
			Text:     line,
			Index:    i,
			FilePath: fileName,
			Selector: "",
		})
	}
	return
}

// getEpubTextLines extracts text lines from an EPUB file or a plain text file
// based on the input file extension and specified start and end positions.
// It returns the extracted text lines, a base name for the output file, and an error if any.
//
// Returns:
// - epubTextLines: A slice of EpubTextLine containing the extracted text lines.
// - outputBaseName: A string representing the base name for the output file.
// - err: An error if any occurred during the extraction process.
func getEpubTextLines() (epubTextLines []epubs.EpubTextLine, outputBaseName string, err error) {
	// epubTextLines = []epubs.EpubTextLine{}
	start := 0
	end := 0

	// check the inpjut file extension
	if path.Ext(config.Params.InputFileName) == ".epub" {
		epubTextLines, start, end, err = epubs.GetEpubTextLines(config.Params.InputFileName, config.Params.Start, config.Params.End, epubs.ProcessableSelectors)
		if err != nil {
			return
		}
	} else {
		// textLines, start, end, err = texts.GetTextFileLines(config.Params.InputFileName, config.Params.Start, config.Params.End)
		epubTextLines, start, end, err = getEpubTextLinesFromTextFile(config.Params.InputFileName, config.Params.Start, config.Params.End)
		if err != nil {
			return
		}
	}
	// calculate a base file name for the output file
	outputBaseName = fmt.Sprintf("%s.lines-%06d-%06d", config.Params.OutputFileName, start, end)
	return
}
