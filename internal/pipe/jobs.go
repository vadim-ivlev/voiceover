// Description: This file contains the job struct and the job queue.

package pipe

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

// JobFunction - function that processes a job.
// It takes a job and returns a modified job and a possible error.
type JobFunction func(Job) (Job, error)

type ProcessLogRecord struct {
	JobID           int           `json:"job_id"`
	Work            string        `json:"work"`
	Worker          string        `json:"worker"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
	Duration        time.Duration `json:"-"`
	DurationSeconds float64       `json:"duration_seconds"`
	Error           string        `json:"error"`
}

type Job struct {
	// Unique identifier for the job
	ID int `json:"id"`
	// Text to process
	Text string `json:"text"`
	// Voice
	Voice string `json:"voice"`
	// Text file
	TextFile string `json:"text_file"`
	// Audio file
	AudioFile string `json:"audio_file"`
	// Process Log. Records of the processing steps for the job
	ProcessLog []ProcessLogRecord `json:"process_log"`
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
