// Description: This file contains the job struct and the job queue.

package pipe

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
)

// JobFunction - function that processes a job.
// It takes a job and returns a modified job and a possible error.
type JobFunction func(Job) (Job, error)

type Task struct {
	Command   string        `json:"command"`
	Params    config.Config `json:"params"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Jobs      []Job         `json:"jobs"`
	Results   struct {
		SoundFile string `json:"sound_file"`
		TextFile  string `json:"text_file"`
	} `json:"results"`
	TaskErrors []string `json:"task_errors"`
}

var pPreviousTask *Task

type WorkLogRecord struct {
	JobID           int           `json:"job_id"`
	Work            string        `json:"work"`
	Worker          string        `json:"worker"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
	Duration        time.Duration `json:"-"`
	DurationSeconds float64       `json:"duration_seconds"`
	Error           string        `json:"error"`
	Result          any           `json:"result"`
}

type Job struct {
	// Unique identifier for the job
	ID      int `json:"id"`
	Results struct {
		// Text to process
		Text string `json:"text"`
		// Voice
		Voice string `json:"voice"`
		// Text file
		TextFile string `json:"text_file"`
		// Audio file
		AudioFile string `json:"audio_file"`
	} `json:"results"`
	// Workers Log. Records of the processing steps for the job
	WorksLog []WorkLogRecord `json:"works_log"`
}

// String - String representation of the job
func (j *Job) String() string {
	return PrettyJSON(j)
}

func CompactJSON(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return string(bytes)
}

func PrettyJSON(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return string(bytes)
}
