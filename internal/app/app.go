// Description: This file contains the main application logic.
package app

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/internal/logger"
)

// InitApp - initializes the application
func InitApp() {
	logger.InitializeLogger()

	// load environment variables from files
	config.ReadConfig(".env")
	config.ReadConfig("voiceover.env")
	// Parse environment variables into the config.Params structure
	config.ParseEnv()

	// if the API key is not set in the config file, then we try to get it from the environment variable
	if config.Params.ApiKey == "" {
		config.Params.ApiKey = os.Getenv("OPENAI_API_KEY")
	}

	RecreateDirs()

	log.Info().Msg("Application started.")

}

// RecreateDirs - recreates directories for texts and sounds
func RecreateDirs() {
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
