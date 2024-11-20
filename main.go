package main

import (
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/pipe"
)

func main() {
	app.InitApp()
	app.ExitIfNoFileToProcess()
	app.RemoveTempFiles()
	log.Info().Msg("Application started.")

	err := pipe.ProcessFile()
	if err != nil {
		log.Error().Msgf("Failed to process file: %v", err)
	}
}
