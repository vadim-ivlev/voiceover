package pipe

import (
	"fmt"
	"time"

	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/pkg/texts"
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

	// Get text lines from the input file
	textLines, start, end, err := texts.GetTextFileLines(config.Params.InputFileName, config.Params.Start, config.Params.End)
	if err != nil {
		return
	}
	// calculate a base file name for the output file
	outputBaseName := fmt.Sprintf("%s.lines-%06d-%06d", config.Params.OutputFileName, start, end)

	// return nil

	// HERE: process the jobs in the pipeline -----------------
	processedJobs, err := DoPipeline(textLines, task)
	if err != nil {
		return
	}

	numDone = len(processedJobs)

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
	if numDone != len(textLines) {
		task.TaskErrors = fmt.Sprintf("Only %d of %d paragraphs processed.", numDone, len(textLines))
	}

	outTaskFile = outputBaseName + ".task.json"
	err = texts.SaveTextFile(outTaskFile, PrettyJSON(task))

	return
}
