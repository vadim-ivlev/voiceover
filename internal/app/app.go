// Description: This file contains the main application logic.
package app

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/pkg/logger"
)

// InitLoggerSetParams - initializes the application
func InitLoggerSetParams() {
	logger.InitializeLogger()
	config.SetAppParams()
	config.PrintAppParams()
}

func ExitIfNoFileToProcess() {
	if config.Params.InputFileName == "" && config.Params.TaskFile == "" {
		log.Error().Msg("No Input File to process.")
		flag.Usage()
		os.Exit(1)
	}
}
