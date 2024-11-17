package sound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
)

type requestBody struct {
	Model string  `json:"model"`
	Input string  `json:"input"`
	Voice string  `json:"voice"`
	Speed float64 `json:"speed"`
}

// alloy, echo, fable, onyx, nova, shimmer
const (
	VoiceAlloy   = "alloy"
	VoiceEcho    = "echo"
	VoiceFable   = "fable"
	VoiceOnyx    = "onyx"
	VoiceNova    = "nova"
	VoiceShimmer = "shimmer"
)

func GenerateMP3(speed float64, voice, text, fileName string) error {
	url := config.Params.BaseURL + "/v1/audio/speech"
	body := requestBody{
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

// GenerateSilenceMP3 - creates a  mp3 file that contains silence using ffmpeg.
// Parameters:
// fileName - the name of the file to create.
// duration - the duration of the silence in seconds.
func GenerateSilenceMP3(fileName string, duration float64) error {
	// Construct the ffmpeg command
	durationStr := fmt.Sprintf("%f", duration)
	log.Info().Msgf("Duration: %s", durationStr)
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "anullsrc=r=44100:cl=mono", "-t", durationStr, "-q:a", "9", "-acodec", "libmp3lame", fileName)

	// Run the ffmpeg command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
