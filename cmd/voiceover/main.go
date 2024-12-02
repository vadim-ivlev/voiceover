package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/pipe"
	"github.com/vadim-ivlev/voiceover/pkg/stopper"
)

func main() {
	startTime := time.Now()

	app.InitLoggerSetParams()
	app.ExitIfNoFileToProcess()

	log.Info().Msg("Application started.")

	// Start watching for the cancel signal
	go stopper.WaitForCancel()

	// Process the input file
	mp3File, txtFile, epubFile, taskFile, numDone, err := pipe.ProcessFile()
	if err != nil {
		log.Error().Msgf("Failed to process file: %v", err)
	} else {
		ResultMessage(mp3File, txtFile, epubFile, taskFile, numDone)
	}

	duration := time.Since(startTime)
	log.Info().Msgf("Operation completed. Duration: %v", duration)
}

func ResultMessage(mp3File, txtFile, epubFile, taskFile string, numDone int) {
	log.Info().Msgf(`
	%d paragraphs processed.
	%d translations done.
	Output files
		MP3 file:  %s
		Text file: %s
		Epub file: %s
		Log file:  %s
	----------------------------	
	`, numDone, pipe.NumTranslations, mp3File, epubFile, txtFile, taskFile)
}
