// Description: This file contains the code to generate speech using the Elevenlabs API.
// Documentation:
// https://elevenlabs.io/docs/api-reference/text-to-speech

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
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style,omitempty"`             // default:"0.0" -greatly improves the speed of generation
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"` // default:"true"
}

type ElevenlabsRequestBody struct {
	Text                            string                  `json:"text"`
	ModelID                         string                  `json:"model_id"`                // default: "eleven_monolingual_v1"
	LanguageCode                    string                  `json:"language_code,omitempty"` //Currently only Turbo v2.5 supports language enforcement.
	VoiceSettings                   ElevenlabsVoiceSettings `json:"voice_settings"`
	PronunciationDictionaryLocators []interface{}           `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            int                     `json:"seed,omitempty"`
	PreviousText                    string                  `json:"previous_text,omitempty"`
	NextText                        string                  `json:"next_text,omitempty"`
	PreviousRequestIDs              []string                `json:"previous_request_ids,omitempty"`
	NextRequestIDs                  []string                `json:"next_request_ids,omitempty"`
	ApplyTextNormalization          string                  `json:"apply_text_normalization,omitempty"` // default=auto  auto, on, off
}

// Elevenlabs voices
const (
	ElevenlabsVoiceSam       string = "ulNeoiyl3bUW7oQjWZE8"
	ElevenlabsVoiceAdamStone string = "NFG5qt843uXKj4pFvR7C"
	ElevenlabsVoiceAlice     string = "Xb7hH8MSUJpSbSDYk0k2"
	ElevenlabsVoiceBrian     string = "7p1URySAeSeJtThZmKB5"
	ElevenlabsVoiceValentino string = "HgJDD5cRFQsVhwzXouaI"
)

func GenerateElevenlabsSpeechMP3(apiKey string, voice, text, fileName string) error {
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voice)
	body := ElevenlabsRequestBody{
		Text:    text,
		ModelID: "eleven_multilingual_v2",
		// LanguageCode: "en", // en-US, en-GB, fr-FR, de-DE, es-ES, it-IT, nl-NL, pt-PT, ru-RU, zh-CN, ja-JP, ko-KR
		VoiceSettings: ElevenlabsVoiceSettings{
			Stability:       0.4,
			SimilarityBoost: 0.8,
			// Style:           0.0, // default: 0.0 -greatly improves the speed of generation
			UseSpeakerBoost: true,
		},
		// PronunciationDictionaryLocators: []interface{}{},
		// Seed:                   123,
		// PreviousText:           "",
		// NextText:               "",
		// PreviousRequestIDs:     []string{},
		// NextRequestIDs:         []string{},
		// ApplyTextNormalization: "auto",
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
	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("xi-api-key", apiKey)

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
