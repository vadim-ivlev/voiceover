package main

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/pipe"
)

func main() {
	app.InitApp()
	log.Info().Msgf("GCLOUD_PROJECT: %s", os.Getenv("GCLOUD_PROJECT"))
	app.ExitIfNoFileToProcess()
	app.RemoveTempFiles()
	log.Info().Msg("Application started.")

	startTime := time.Now()
	mp3File, txtFile, lofFile, err := pipe.ProcessFile()
	if err != nil {
		log.Error().Msgf("Failed to process file: %v", err)
	} else {
		log.Info().Msg("File processed successfully.")
		log.Info().Msgf("MP3 file: %s", mp3File)
		log.Info().Msgf("Text file: %s", txtFile)
		log.Info().Msgf("Log file: %s", lofFile)
	}
	// Log duration of the operation
	duration := time.Since(startTime)
	log.Info().Msgf("Operation completed. Duration: %v", duration)
}
