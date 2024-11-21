package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v11"
)

// config - структура для хранения параметров приложения
type config struct {
	// env
	BaseURL          string `env:"BASE_URL" envDefault:"https://api.openai.com"`
	ApiKey           string `env:"API_KEY"`
	TextsDir         string `env:"TEXTS_DIR" envDefault:"./.data/texts"`
	SoundsDir        string `env:"SOUNDS_DIR" envDefault:"./.data/sounds"`
	FileListFileName string `env:"FILE_LIST_FILE_NAME" envDefault:"file-list.txt"`
	// command line
	Start          int    `json:"-s: Start      "`
	End            int    `json:"-e: End        "`
	Voices         string `json:"-v: Voices     "`
	OutputFileName string `json:"-o: Output File"`
	InputFileName  string `json:"  : Input  File"`
}

func (c config) String() string {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return string(bytes)
}

// Params - переменная для хранения параметров приложения
var Params config = config{}

func customUsage() {
	fmt.Fprintf(os.Stderr, `
This application processes a text file and generates voice output.

Usage:
    >>> %s [OPTIONS] <input file name>
Example of usage:
    >>> %s "Story.txt"

`, os.Args[0], os.Args[0])
	fmt.Fprintf(os.Stderr, "List of complementary OPTIONS:\n")
	flag.PrintDefaults()
}

// ParseCommandLine - читает параметры командной строки с значениями по умолчанию
func ParseCommandLine() {
	flag.Usage = customUsage

	flag.IntVar(&Params.Start, "s", 0, "Number of the first line of the file to process. Starting from 0.")
	flag.IntVar(&Params.End, "e", 0, "Number of the last line of the file to process. Last line will not be processed. 0 - process to the end of the file.")
	flag.StringVar(&Params.Voices, "v", "alloy,echo", "Comma separated list of voices (alloy,echo,fable,onyx,nova,shimmer). No spaces.")

	flag.StringVar(&Params.OutputFileName, "o", "", "Output file name. If empty, will be equal to the input file name.")
	flag.Parse()

	Params.InputFileName = flag.Arg(0)
	if Params.OutputFileName == "" {
		Params.OutputFileName = Params.InputFileName
	}
}

// ParseEnv - читает переменные окружения (если они есть) и сохраняет их в структуру Params
func ParseEnv() {
	if err := env.Parse(&Params); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

// ReadConfig reads env file and fill EnvParams with environment variables values
func ReadConfig(fileName string) {
	if err := godotenv.Load(fileName); err != nil {
		log.Warn().Msg(err.Error())
	}
}

// JSONString - сериализуем структуру в формат JSON
func JSONString(params interface{}) string {
	bytes, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return string(bytes)
}

// SetAppParams - устанавливает параметры приложения
func SetAppParams() {
	// load environment variables from files
	ReadConfig(".env")
	ReadConfig("voiceover.env")
	// Parse environment variables into the Params structure
	ParseEnv()

	// if the API key is not set in the config file, then we try to get it from the environment variable
	if Params.ApiKey == "" {
		Params.ApiKey = os.Getenv("OPENAI_API_KEY")
	}

	// Parse command line arguments
	ParseCommandLine()
}

// PrintAppParams - выводит параметры приложения
func PrintAppParams() {
	// Print config
	log.Info().Msg(JSONString(Params))
}
