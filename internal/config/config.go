package config

import (
	"encoding/json"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v11"
)

// config - структура для хранения параметров приложения
type config struct {
	BaseURL   string `env:"BASE_URL" envDefault:"https://api.openai.com"`
	ApiKey    string `env:"API_KEY"`
	TextsDir  string `env:"TEXTS_DIR" envDefault:"./.data/texts"`
	SoundsDir string `env:"SOUNDS_DIR" envDefault:"./.data/sounds"`
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

// // ParseCommandLine - читает параметры командной строки с значениями по умолчанию
// func ParseCommandLine() {
// 	flag.StringVar(&Params.ServerAddress, "a", "localhost:8080", "HTTP server address")
// 	flag.StringVar(&Params.BaseURL, "b", "http://localhost:8080", "Base URL")
// 	flag.StringVar(&Params.FileStoragePath, "f", "./data/file-storage.txt", "File storage path")
// 	flag.StringVar(&Params.DatabaseDSN, "d", "", "Database DSN")
// 	flag.Parse()
// }

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
