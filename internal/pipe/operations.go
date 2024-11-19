package pipe

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
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

// StartPipeline - starts the pipeline processing.
// Parameters:
// textLines: the text lines to process
// Returns:
// an array of processed jobs
// an error if any
func StartPipeline(textLines []string) (doneJobs []Job, err error) {

	// Create job channels
	jobsChan := make(chan Job)
	textChan := make(chan Job, 50)
	soundChan := make(chan Job, 50)

	// go GenerateJobs(10, 0, jobsChan)
	go toChannel(newJobsArray(10), jobsChan, 0)
	doTeamWork(1, "T", getTextOperation(textLines), jobsChan, textChan)
	doTeamWork(1, "S", soundOperation, textChan, soundChan)

	// gatther the jobs into an array
	doneJobs = toArray(soundChan)
	return doneJobs, nil
}

func ProcessFile(filePath string) (err error) {
	// split the file
	textLines, err := text.SplitTextFileScan(filePath)
	if err != nil {
		return err
	}
	log.Info().Msg(PrettyJSON(textLines))

	// process the jobs in the pipeline
	processedJobs, err := StartPipeline(textLines)
	if err != nil {
		return err
	}

	// log the processed jobs
	log.Info().Msg("All jobs have been processed.<<<<<<<<<<<<<<<<<<<")
	log.Info().Msg(PrettyJSON(processedJobs))
	return nil
}
