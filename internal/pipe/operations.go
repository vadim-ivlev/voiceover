package pipe

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/internal/sound"
	"github.com/vadim-ivlev/voiceover/internal/text"
)

// createJob - create a job with the given id
func createJob(id int) Job {
	job := Job{
		ID: id,
	}
	return job
}

// LogJob - logs the job information
func LogJob(job Job, msg string) {
	log.Info().Msgf("Job %d %s", job.ID, msg)
	// log.Info().Msg(job.String())
}

// getTextJobOperation - returns the textJob operation
func getTextOperation(textLines []string) JobFunction {
	// adds text to the job, and assigns a voice
	return func(job Job) Job {
		if len(textLines) < job.ID+1 {
			log.Error().Msgf("Job %d: No text for the job", job.ID)
			return job
		}
		job.Text = strings.Trim(textLines[job.ID], " \r\n")
		if len(job.Text) > 0 {
			job.Voice = sound.NextVoice()
		}
		return job
	}
}

// soundOperation - generates a sound file and a text file for the job
func soundOperation(job Job) Job {
	job.TextFile = fmt.Sprintf("%s/%08d.txt", config.Params.TextsDir, job.ID)
	job.AudioFile = fmt.Sprintf("%s/%08d.mp3", config.Params.SoundsDir, job.ID)

	err := sound.GenerateMP3(0.5, 1.0, job.Voice, job.Text, job.AudioFile)
	if err != nil {
		log.Error().Msgf("Job %d: Failed to generate sound file: %v", job.ID, err)
		job.RequestError = err.Error()
	}

	err = text.SaveTextFile(job.TextFile, job.Text)
	if err != nil {
		log.Error().Msgf("Job %d: Failed to save text file: %v", job.ID, err)
		job.RequestError += "; " + err.Error()
	}

	return job
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
func DoPipeline(textLines []string) (doneJobs []Job, err error) {

	numLines := len(textLines)

	// Create channels to pass the jobs between the workers
	jobsChan := make(chan Job)
	textChan := make(chan Job)
	// soundChan must be buffered because its jobs will be sorted after it is closed
	soundChan := make(chan Job, numLines)

	// Fill the jobs channel with the jobs from a newly created array
	go toChannel(newJobsArray(numLines), jobsChan, 0)
	// add a text to each job and assign a voice
	// doTeamWork(1, "T", getTextOperation(textLines), jobsChan, textChan)
	go DoWork(nil, "T", getTextOperation(textLines), jobsChan, textChan)
	// generate sound file for each job. Fan-out.
	doTeamWork(10, "S", soundOperation, textChan, soundChan)

	// gatther the jobs into an array. Fan-in.
	doneJobs = toArray(soundChan)
	return doneJobs, nil
}

// ProcessFile - processes the input file.
func ProcessFile() (outMP3File string, outTextFile string, err error) {

	// Get text lines from the input file
	textLines, start, end, err := text.GetTextFileLines(config.Params.InputFileName, config.Params.Start, config.Params.End)
	if err != nil {
		return
	}

	// // Print extracted text lines
	// for i, line := range textLines {
	// 	log.Info().Msgf("%06d: %s", start+i, line)
	// }

	// calculate file name for the output file
	outputFileName := fmt.Sprintf("%s.lines-%06d-%06d", config.Params.OutputFileName, start, end)

	//remove the output mp3 file if it exists
	err = os.Remove(outputFileName + ".mp3")
	if err != nil {
		log.Info().Msgf("Failed to delete the output file: %v", err)
	}

	// remove the output text file if it exists
	err = os.Remove(outputFileName + ".txt")
	if err != nil {
		log.Info().Msgf("Failed to delete the output file: %v", err)
	}

	// clear directory of text and sound files
	app.RemoveTempFiles()

	// return nil

	// HERE: process the jobs in the pipeline -----------------
	processedJobs, err := DoPipeline(textLines)
	if err != nil {
		return
	}

	// create file list of audio files for concatenation
	err = CreateFileList(processedJobs)
	if err != nil {
		return
	}

	// concatenate the audio files into one
	err = sound.ConcatenateMP3Files(config.Params.FileListFileName, outputFileName+".mp3")
	if err != nil {
		return
	}
	outMP3File = outputFileName + ".mp3"

	// write a text file with processed lines
	err = text.SaveTextFile(outputFileName+".txt", strings.Join(textLines, "\n"))
	if err != nil {
		return
	}
	outTextFile = outputFileName + ".txt"

	// write log of processed jobs
	err = text.SaveTextFile(outputFileName+".log.json", PrettyJSON(processedJobs))
	if err != nil {
		return
	}

	return
}

// CreateFileList - creates a file list of audio files for concatenation with ffmpeg
// File list example:
// file 'file1.mp3'
// file 'file2.mp3'
// file 'file3.mp3'
// Parameters:
// jobs: the jobs to process
func CreateFileList(jobs []Job) (err error) {
	fileList := ""
	for _, job := range jobs {
		fileList += fmt.Sprintf("file '%s'\n", job.AudioFile)
	}
	err = text.SaveTextFile(config.Params.FileListFileName, fileList)
	if err != nil {
		return err
	}
	return nil
}
