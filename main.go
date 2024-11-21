package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/pipe"
)

func main() {
	app.InitApp()
	app.ExitIfNoFileToProcess()
	app.RemoveTempFiles()
	log.Info().Msg("Application started.")

	startTime := time.Now()
	mp3File, txtFile, err := pipe.ProcessFile()
	if err != nil {
		log.Error().Msgf("Failed to process file: %v", err)
	} else {
		log.Info().Msg("File processed successfully.")
		log.Info().Msgf("MP3 file: %s", mp3File)
		log.Info().Msgf("Text file: %s", txtFile)
	}
	duration := time.Since(startTime)
	log.Info().Msgf("Duration: %v", duration)

}
