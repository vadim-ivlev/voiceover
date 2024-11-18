// Description: This file contains the job struct and the job queue.

package pipe

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

type ProcessLogRecord struct {
	JobID       int           `json:"job_id"`
	WorkerName  string        `json:"worker_name"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	DurationSec float64       `json:"duration_sec"`
	Error       string        `json:"error"`
}

type Job struct {
	// Unique identifier for the job
	ID int `json:"id"`
	// Text to process
	Text string `json:"text,omitempty"`
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
	// // Processed by worker
	// ProcessedBy string `json:"processed_by,omitempty"`
	// // Result Messages
	// ResultMessages string `json:"result_messages,omitempty"`
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
