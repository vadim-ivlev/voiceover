package sound

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
)

type googleTTSRequest struct {
	Input struct {
		Text string `json:"text"`
	} `json:"input"`
	Voice struct {
		LanguageCode string `json:"languageCode"`
		Name         string `json:"name"`
	} `json:"voice"`
	AudioConfig struct {
		AudioEncoding string  `json:"audioEncoding"`
		SpeakingRate  float64 `json:"speakingRate,omitempty"`
	} `json:"audioConfig"`
}

type googleTTSResponse struct {
	AudioContent string `json:"audioContent"`
}

func GenerateGoogleSpeechMP3(speed float64, voice, text, fileName string) error {
	url := "https://texttospeech.googleapis.com/v1/text:synthesize"
	_ = voice
	voice = "en-GB-Journey-D"

	reqBody := googleTTSRequest{}
	reqBody.Input.Text = text
	reqBody.Voice.LanguageCode = "en-GB"
	reqBody.Voice.Name = voice
	reqBody.AudioConfig.AudioEncoding = "LINEAR16" //LINEAR16 is for WAV, MP3 is for MP3
	if speed != 1.0 {
		reqBody.AudioConfig.SpeakingRate = speed
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	prettyJson, _ := json.MarshalIndent(reqBody, "", "  ")
	log.Info().Msgf("Request: %s", string(prettyJson))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	log.Info().Msgf("GCLOUD_PROJECT: %s", os.Getenv("GCLOUD_PROJECT"))
	log.Info().Msgf("GCLOUD_ACCESS_TOKEN: %s", os.Getenv("GCLOUD_ACCESS_TOKEN"))

	req.Header.Set("X-Goog-User-Project", config.Params.GcloudProject)
	req.Header.Set("Authorization", "Bearer "+config.Params.GcloudAccessToken)
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", "Bearer "+config.Params.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.Status)
	}

	var ttsResp googleTTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&ttsResp); err != nil {
		return err
	}

	audio, err := base64.StdEncoding.DecodeString(ttsResp.AudioContent)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileName, audio, 0644); err != nil {
		return err
	}

	return nil
}
