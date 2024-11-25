package pipe

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/internal/text"
)

// ********************************************************************************************************************

// ProcessFile - processes the input file.
func ProcessFile() (outMP3File, outTextFile, outTaskFile string, err error) {

	startTime := time.Now()

	var task = Task{
		Params: config.Params,
	}

	// check if we should continue a previous task
	if config.Params.TaskFile != "" {
		err = LoadJSONFile(config.Params.TaskFile, pPreviousTask)
		if err != nil {
			return
		}
		task = *pPreviousTask
	}

	task.StartTime = startTime
	task.Command = strings.Join(os.Args, " ")

	// Get text lines from the input file
	textLines, start, end, err := text.GetTextFileLines(config.Params.InputFileName, config.Params.Start, config.Params.End)
	if err != nil {
		return
	}

	// calculate a base file name for the output file
	outputBaseName := fmt.Sprintf("%s.lines-%06d-%06d", config.Params.OutputFileName, start, end)

	// remove the output mp3 file if it exists
	// TODO: move down
	err = os.Remove(outputBaseName + ".mp3")
	if err != nil {
		log.Info().Msgf("Failed to delete the output file: %v", err)
	}

	// remove the output text file if it exists
	err = os.Remove(outputBaseName + ".txt")
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

	// join the processed jobs into one mp3 file
	outMP3File, err = JoinMP3Files(processedJobs, outputBaseName)
	if err != nil {
		return
	}

	// write a text file with processed lines
	outTextFile = outputBaseName + ".txt"
	err = text.SaveTextFile(outTextFile, strings.Join(textLines, "\n"))
	if err != nil {
		return
	}

	endTime := time.Now()
	task.EndTime = endTime
	task.Duration = endTime.Sub(startTime)
	task.Jobs = processedJobs

	outTaskFile = outputBaseName + ".task.json"
	err = text.SaveTextFile(outTaskFile, PrettyJSON(task))

	return
}
