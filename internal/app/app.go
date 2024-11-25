// Description: This file contains the main application logic.
package app

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/internal/logger"
)

// InitApp - initializes the application
func InitApp() {
	logger.InitializeLogger()
	config.SetAppParams()
	config.PrintAppParams()
}

// RemoveTempFiles - recreates directories for texts and sounds
func RemoveTempFiles() {
	err := os.RemoveAll(config.Params.FileListFileName)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
	err = os.RemoveAll(config.Params.TextsDir)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
	err = os.RemoveAll(config.Params.SoundsDir)
	if err != nil {
		log.Warn().Msg(err.Error())
	}
	//-----------------

	err = os.MkdirAll(config.Params.TextsDir, 0755)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	err = os.MkdirAll(config.Params.SoundsDir, 0755)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func ExitIfNoFileToProcess() {
	if config.Params.InputFileName == "" {
		log.Error().Msg("No Input File to process.")
		flag.Usage()
		os.Exit(1)
	}
}
