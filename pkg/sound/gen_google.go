package sound

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	voice = "en-GB-Journey-D" // en-GB-News-K // en-GB-Studio-B M // en-GB-Studio-C F

	ttsReq := googleTTSRequest{}
	ttsReq.Input.Text = text
	ttsReq.Voice.LanguageCode = "en-GB"
	ttsReq.Voice.Name = voice
	ttsReq.AudioConfig.AudioEncoding = "MP3" //LINEAR16 is for WAV, MP3 is for MP3
	if speed != 1.0 {
		ttsReq.AudioConfig.SpeakingRate = speed
	}

	ttsReqJSON, err := json.Marshal(ttsReq)
	if err != nil {
		return err
	}

	// prettyJson, _ := json.MarshalIndent(reqBody, "", "  ")
	// log.Info().Msgf("Request: %s", string(prettyJson))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(ttsReqJSON))
	if err != nil {
		return err
	}

	req.Header.Set("X-Goog-User-Project", os.Getenv("GCLOUD_PROJECT"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GCLOUD_ACCESS_TOKEN"))
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
