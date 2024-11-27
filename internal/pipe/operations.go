package pipe

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/audio"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/pkg/sound"
	"github.com/vadim-ivlev/voiceover/pkg/text"
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

// getTextJobOperation - returns the textJob operation
func getTextOperation(textLines []string) JobFunction {
	// adds text to the job, and assigns a voice
	return func(job Job) (Job, error) {
		if len(textLines) < job.ID+1 {
			return job, fmt.Errorf("Job %d: No text for the job", job.ID)
		}
		job.Results.Text = strings.TrimSpace(textLines[job.ID])
		return job, nil
	}
}

// soundOperation - generates a sound file and a text file for the job
func soundOperation(job Job) (Job, error) {
	job.Results.TextFile = fmt.Sprintf("%s/%08d.txt", config.Params.TextsDir, job.ID)
	job.Results.AudioFile = fmt.Sprintf("%s/%08d.mp3", config.Params.SoundsDir, job.ID)

	// select voice
	if len(job.Results.Text) > 0 {
		job.Results.Voice = audio.NextVoice()
	}

	var err error

	// if voice is empty, generate silence
	if job.Results.Voice == "" {

		err = sound.GenerateSilenceMP3(config.Params.Pause, job.Results.AudioFile)
	} else {
		err = audio.GenerateSpeechMP3(config.Params.Speed, job.Results.Voice, job.Results.Text, job.Results.AudioFile)
	}
	if err != nil {
		return job, err
	}

	err = text.SaveTextFile(job.Results.TextFile, job.Results.Text)
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
// textLines: the text lines to process
// Returns:
// an array of processed jobs
// an error if any
func DoPipeline(textLines []string, task Task) (doneJobs []Job, err error) {

	numLines := len(textLines)

	jobs := newJobsArray(numLines)
	loadPreviuosJobs(jobs, task)

	// Create channels to pass the jobs between the workers
	jobsChan := make(chan Job)
	textChan := make(chan Job)
	// soundChan must be buffered because its jobs will be sorted after it is closed
	soundChan := make(chan Job, numLines)

	// Fill the jobs channel with the jobs from a newly created array
	go toChannel(jobs, jobsChan)
	// add a text to each job and assign a voice
	// doTeamWork(1, "T", getTextOperation(textLines), jobsChan, textChan)
	go DoWork(nil, "Text", "T", getTextOperation(textLines), jobsChan, textChan)
	// generate sound file for each job. Fan-out.
	doTeamWork(10, "Sound", "S", soundOperation, textChan, soundChan)

	// gatther the jobs into an array. Fan-in.
	doneJobs = toArray(soundChan)
	return doneJobs, nil
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
	err = text.SaveTextFile(config.Params.FileListFileName, fileList)
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
		textLines = append(textLines, job.Results.Text)
	}

	// write a text file with processed lines
	outTextFile = outputBaseName + ".txt"
	err = text.SaveTextFile(outTextFile, strings.Join(textLines, "\n"))
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
