package sound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ElevenlabsVoiceSettings struct {
	Stability       int  `json:"stability"`
	SimilarityBoost int  `json:"similarity_boost"`
	Style           int  `json:"style"`
	UseSpeakerBoost bool `json:"use_speaker_boost"`
}

type ElevenlabsRequestBody struct {
	Text                            string                  `json:"text"`
	ModelID                         string                  `json:"model_id"`
	LanguageCode                    string                  `json:"language_code"`
	VoiceSettings                   ElevenlabsVoiceSettings `json:"voice_settings"`
	PronunciationDictionaryLocators []interface{}           `json:"pronunciation_dictionary_locators"`
	Seed                            int                     `json:"seed"`
	PreviousText                    string                  `json:"previous_text"`
	NextText                        string                  `json:"next_text"`
	PreviousRequestIDs              []string                `json:"previous_request_ids"`
	NextRequestIDs                  []string                `json:"next_request_ids"`
	UsePvcAsIvc                     bool                    `json:"use_pvc_as_ivc"`
	ApplyTextNormalization          string                  `json:"apply_text_normalization"`
}

func GenerateElevenlabsSpeechMP3(speed float64, voice, text, fileName string) error {
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voice)
	body := ElevenlabsRequestBody{
		Text:         text,
		ModelID:      "your_model_id",
		LanguageCode: "your_language_code",
		VoiceSettings: ElevenlabsVoiceSettings{
			Stability:       123,
			SimilarityBoost: 123,
			Style:           123,
			UseSpeakerBoost: true,
		},
		PronunciationDictionaryLocators: []interface{}{},
		Seed:                            123,
		PreviousText:                    "",
		NextText:                        "",
		PreviousRequestIDs:              []string{},
		NextRequestIDs:                  []string{},
		UsePvcAsIvc:                     true,
		ApplyTextNormalization:          "auto",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to generate speech: %s", resp.Status)
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
