package pipe

import (
	"fmt"
	"time"

	"github.com/vadim-ivlev/voiceover/pkg/texts"
	"github.com/vadim-ivlev/voiceover/pkg/utils"
)

// ********************************************************************************************************************

// ProcessFile - processes the input file.
func ProcessFile() (outMP3File, outTextFile, outTaskFile string, numDone int, err error) {
	startTime := time.Now()

	// Create a new task or restore the previous one
	task, err := CreateOrRestoreTask()
	if err != nil {
		return
	}

	// PROCESS PROCESS PROCESS PROCESS the jobs in the pipeline -----------------
	processedJobs, numJobs, outputBaseName, err := DoPipeline(task)
	if err != nil {
		return
	}

	outMP3File, err = CreateOutputMP3(processedJobs, outputBaseName)
	if err != nil {
		return
	}

	outTextFile, err = CreateOutputText(processedJobs, outputBaseName)
	if err != nil {
		return
	}

	task.StartTime = startTime
	task.EndTime = time.Now()
	task.Duration = task.EndTime.Sub(task.StartTime)
	task.Jobs = processedJobs
	task.Results.SoundFile = outMP3File
	task.Results.TextFile = outTextFile
	if len(processedJobs) != numJobs {
		task.TaskErrors = fmt.Sprintf("Only %d of %d jobs processed.", len(processedJobs), numJobs)
	}

	outTaskFile = outputBaseName + ".task.json"
	err = texts.SaveTextFile(outTaskFile, utils.PrettyJSON(task))

	return
}
