// Description: This file contains the job struct and the job queue.

package pipe

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

// JobFunction - function that processes a job.
// It takes a job and returns a modified job.
type JobFunction func(Job) Job

type ProcessLogRecord struct {
	JobID           int           `json:"job_id"`
	WorkerName      string        `json:"worker_name"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"-"`
	Duration        time.Duration `json:"-"`
	DurationSeconds float64       `json:"duration_seconds"`
	Error           string        `json:"-"`
}

type Job struct {
	// Unique identifier for the job
	ID int `json:"id"`
	// Text to process
	Text string `json:"text"`
	// Voice
	Voice string `json:"voice,omitempty"`
	// Request time
	RequestTime string `json:"request_time,omitempty"`
	// Request duration
	RequestDuration string `json:"request_duration,omitempty"`
	// Request error
	RequestError string `json:"request_error,omitempty"`
	// Text file
	TextFile string `json:"text_file,omitempty"`
	// Audio file
	AudioFile string `json:"audio_file,omitempty"`
	// Process Log. Records of the processing steps for the job
	ProcessLog []ProcessLogRecord `json:"process_log,omitempty"`
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

// String - String representation of the job
func (j *Job) String() string {
	return PrettyJSON(j)
}
