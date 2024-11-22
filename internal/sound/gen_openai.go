package sound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/vadim-ivlev/voiceover/internal/config"
)

type openaiRequestBody struct {
	Model string  `json:"model"`
	Input string  `json:"input"`
	Voice string  `json:"voice"`
	Speed float64 `json:"speed"`
}

// GenerateOpenaiSpeechMP3 generates an MP3 file with the given text using the OpenAI API.
func GenerateOpenaiSpeechMP3(speed float64, voice, text, fileName string) error {
	url := config.Params.BaseURL + "/v1/audio/speech"
	body := openaiRequestBody{
		Model: "tts-1",
		Input: text,
		Voice: voice,
		Speed: speed,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.Params.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.Status)
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}
