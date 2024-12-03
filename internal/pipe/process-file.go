package pipe

import (
	"fmt"
	"sync"
	"time"

	"github.com/vadim-ivlev/voiceover/pkg/epubs"
	"github.com/vadim-ivlev/voiceover/pkg/texts"
	"github.com/vadim-ivlev/voiceover/pkg/utils"
)

// ********************************************************************************************************************

// ProcessFile - processes the input file.
func ProcessFile() (outMP3File, outTextFile, outEpubFile, outTaskFile string, numDone int, err error) {
	startTime := time.Now()

	// Create a new task or restore the previous one
	task, err := CreateOrRestoreTask()
	if err != nil {
		return
	}

	// Wait group for the toc translation
	wgTableOfContents := sync.WaitGroup{}
	wgTableOfContents.Add(1)
	go func() {
		_, err = epubs.TranslateTablesOfContents(task.Params.InputFileName)
		if err != nil {
			task.TaskErrors = fmt.Sprintf(" !!!!!!!!!! Error translating TOC: %v", err)
		}
		wgTableOfContents.Done()
	}()

	// PROCESS PROCESS PROCESS PROCESS the jobs in the pipeline -----------------
	processedJobs, numJobs, outputBaseName, err := DoPipeline(task)
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

	// Wait for the TOC translation to finish before creating the output EPUB
	wgTableOfContents.Wait()

	outEpubFile, err = CreateOutputEpub(task.Params.InputFileName, processedJobs, outputBaseName)
	if err != nil {
		return
	}

	task.StartTime = startTime
	task.EndTime = time.Now()
	task.Duration = task.EndTime.Sub(task.StartTime)
	task.Jobs = processedJobs
	task.Results.SoundFile = outMP3File
	task.Results.TextFile = outTextFile
	task.Results.EpubFile = outEpubFile
	if len(processedJobs) != numJobs {
		task.TaskErrors = fmt.Sprintf("Only %d of %d jobs processed.", len(processedJobs), numJobs)
	}

	outTaskFile = outputBaseName + ".task.json"
	err = texts.SaveTextFile(outTaskFile, utils.PrettyJSON(task))

	return
}
