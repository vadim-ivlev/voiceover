package main

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/pipe"
	"github.com/vadim-ivlev/voiceover/internal/stopper"
)

func main() {
	app.InitApp()
	log.Info().Msgf("GCLOUD_PROJECT: %s", os.Getenv("GCLOUD_PROJECT"))
	app.ExitIfNoFileToProcess()
	app.RemoveTempFiles()
	log.Info().Msg("Application started.")

	// Start watching for the cancel signal
	go stopper.WaitForCancel()

	startTime := time.Now()
	mp3File, txtFile, taskFile, err := pipe.ProcessFile()
	if err != nil {
		log.Error().Msgf("Failed to process file: %v", err)
	} else {
		log.Info().Msg("File processed successfully.")
		log.Info().Msgf("MP3 file: %s", mp3File)
		log.Info().Msgf("Text file: %s", txtFile)
		log.Info().Msgf("Log file: %s", taskFile)
	}
	// Log duration of the operation
	duration := time.Since(startTime)
	log.Info().Msgf("Operation completed. Duration: %v", duration)
}
