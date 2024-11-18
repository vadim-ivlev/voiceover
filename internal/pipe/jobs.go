// Description: This file contains the job struct and the job queue.

package pipe

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

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
	// Processed by worker
	ProcessedBy string `json:"processed_by,omitempty"`
	// Result Messages
	ResultMessages string `json:"result_messages,omitempty"`
}

// String - String representation of the job
func (j *Job) String() string {
	bytes, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return string(bytes)
}
