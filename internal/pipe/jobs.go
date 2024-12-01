// Description: This file contains the job struct and the job queue.

package pipe

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/pkg/epubs"
	"github.com/vadim-ivlev/voiceover/pkg/utils"
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
	TaskErrors string `json:"task_errors"`
}

// var pPreviousTask *Task

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
		// EPUB properties
		Epub epubs.EpubTextLine `json:"epub"`
		// Text to process
		Text string `json:"text"`
		// Translated text
		TranslatedText string `json:"translated_text"`
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
	return utils.PrettyJSON(j)
}

func RemoveTempDirs() {
	err := os.RemoveAll(config.Params.TextsDir)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
	err = os.RemoveAll(config.Params.SoundsDir)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
}

func CreateTempDirs() {
	err := os.MkdirAll(config.Params.TextsDir, 0755)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	err = os.MkdirAll(config.Params.SoundsDir, 0755)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

// CreateOrRestoreTask - Create a new task or restore the previous one/
// If a previous task file is provided, the task parameters are copied from the previous task.
func CreateOrRestoreTask() (task Task, err error) {
	task = Task{}

	// check if we should continue a previous task
	if config.Params.TaskFile == "" {
		// clear directory of text and sound files
		RemoveFileListFile()
		RemoveTempDirs()
		CreateTempDirs()
		task.Command = strings.Join(os.Args, " ")
		task.Params = config.Params
	} else {
		pPreviousTask := new(Task)
		err = LoadJSONFile(config.Params.TaskFile, pPreviousTask)
		if err != nil {
			return
		}
		task = *pPreviousTask
		// copy the previous task parameters to the current task
		config.Params = task.Params
		CreateTempDirs()
	}
	return
}
