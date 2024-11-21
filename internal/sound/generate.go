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

var currentVoice = -1

// Returns next voice in the circular list of voices.
func NextVoice() string {
	voices := []string{VoiceAlloy, VoiceEcho, VoiceFable, VoiceOnyx, VoiceNova, VoiceShimmer}
	currentVoice++
	if currentVoice >= len(voices) || currentVoice < 0 {
		currentVoice = 0
	}
	return voices[currentVoice]
}

func GenerateSpeechMP3(speed float64, voice, text, fileName string) error {
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
func GenerateSilenceMP3(duration float64, fileName string) error {
	// Construct the ffmpeg command
	durationStr := fmt.Sprintf("%f", duration)
	// log.Info().Msgf("Duration: %s", durationStr)
	// cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "anullsrc=r=44100:cl=mono", "-t", durationStr, "-q:a", "9", "-acodec", "libmp3lame", fileName)
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "anullsrc=r=24000:cl=mono", "-t", durationStr, "-q:a", "9", "-acodec", "libmp3lame", fileName)

	// Run the ffmpeg command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// GenerateMP3 - creates a mp3 file that contains silence or speech.
// Parameters:
// silenceDuration - the duration of the silence in seconds if text is empty.
// speed - the speed of the speech.
// voice - the voice to use for the speech.
// text - the text to speak.
// fileName - the name of the file to create.
func GenerateMP3(silenceDuration, speed float64, voice, text, fileName string) (err error) {
	if len(text) > 0 {
		return GenerateSpeechMP3(speed, voice, text+".", fileName)
	} else {
		return GenerateSilenceMP3(silenceDuration, fileName)
	}
}

/*
To combine MP3 files using FFmpeg, follow these steps:

1. Create a text file listing the MP3 files to merge
Create a text file, e.g., file_list.txt, with the following format:

file 'file1.mp3'
file 'file2.mp3'
file 'file3.mp3'

Each line should start with file followed by the file name enclosed in single quotes. Ensure the file paths are correct.

2. Run the FFmpeg command
Use the following command to merge the MP3 files:

bash

ffmpeg -f concat -safe 0 -i file_list.txt -c copy output.mp3
*/

// ConcatenateMP3Files - concatenates a list of mp3 files into a single mp3 file using ffmpeg.
// Parameters:
// fileListFileName - the name of the file that contains the list of mp3 files to concatenate.
// outputFileName - the name of the output file.
func ConcatenateMP3Files(fileListFileName, outputFileName string) error {
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", fileListFileName, "-c:a", "libmp3lame", "-q:a", "2", outputFileName, "-y")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Error().Msg("*****************************************************")
		log.Printf("Failed to concatenate mp3 files: %v\nStderr: %s", err, stderr.String())
		log.Error().Msg("*****************************************************")
		return err
	}
	return nil
}
