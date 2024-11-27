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

// Config - структура для хранения параметров приложения
type Config struct {
	// env
	GcloudProject     string `json:"gcloud_project" env:"GCLOUD_PROJECT"`
	GcloudAccessToken string `json:"gcloud_access_token" env:"GCLOUD_ACCESS_TOKEN"`
	ElevenlabsAPIKey  string `json:"elevenlabs_api_key" env:"ELEVENLABS_API_KEY"`
	OpenaiAPIURL      string `json:"openai_api_url" env:"OPENAI_API_URL" envDefault:"https://api.openai.com"`
	ApiKey            string `json:"api_key" env:"API_KEY"`

	// command line
	TextsDir         string  `json:"texts_dir"`
	SoundsDir        string  `json:"sounds_dir"`
	FileListFileName string  `json:"file_list_file_name"`
	TTSAPI           string  `json:"tts_api"`
	Start            int     `json:"start"`
	End              int     `json:"end"`
	Voices           string  `json:"voices"`
	OutputFileName   string  `json:"output_file_name"`
	InputFileName    string  `json:"input_file_name"`
	TaskFile         string  `json:"task_file"`
	Speed            float64 `json:"speed"`
	Pause            float64 `json:"pause"`
	TranslateTo      string  `json:"translate_to"`

	// debug variables
	NapTime int `json:"nap_time"`
}

func (c Config) String() string {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return string(bytes)
}

// Params - переменная для хранения параметров приложения
var Params Config = Config{}

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

	flag.StringVar(&Params.TextsDir, "textsdir", "./.data/texts", "Directory with text files.")
	flag.StringVar(&Params.SoundsDir, "soundsdir", "./.data/sounds", "Directory with sound files.")
	flag.StringVar(&Params.FileListFileName, "filelist", "file-list.txt", "File with a list of MP3  files to concatenate.")

	flag.StringVar(&Params.TTSAPI, "ttsapi", "openai", "TTS API to use {openai|google|elevenlabs}.")
	flag.IntVar(&Params.Start, "s", 0, "Number of the first line of the file to process. Starting from 0.")
	flag.IntVar(&Params.End, "e", 0, "Number of the last line of the file to process. Last line will not be processed. 0 - process to the end of the file.")
	flag.StringVar(&Params.Voices, "voices", "", "Comma separated list of voices (alloy,echo,fable,onyx,nova,shimmer). No spaces.")
	flag.StringVar(&Params.OutputFileName, "o", "", "Output file name. If empty, will be equal to the input file name.")
	flag.StringVar(&Params.TaskFile, "task", "", "Previous task file to continue processing.")
	flag.Float64Var(&Params.Speed, "speed", 1.0, "Speed of the voice.")
	flag.Float64Var(&Params.Pause, "pause", 0.7, "Pause between paragraphs in seconds.")
	flag.StringVar(&Params.TranslateTo, "translateto", "", "Translate text to the given language. Russian, German, etc.")

	flag.IntVar(&Params.NapTime, "nap", 0, "Random nap time up to the given value in milliseconds between worker operations")
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

// ReadEnvFile reads env file and fill EnvParams with environment variables values
func ReadEnvFile(fileName string) {
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
	ReadEnvFile(".env")
	ReadEnvFile("voiceover.env")
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
