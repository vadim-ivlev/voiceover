package audio

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/pkg/sound"
)

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
	// get voices from config
	voices := strings.Split(config.Params.Voices, ",")
	if len(voices) == 0 {
		return ""
	}
	currentVoice = (currentVoice + 1) % len(voices)
	return voices[currentVoice]
}

// GenerateSpeechMP3 - creates a mp3 file that contains speech.
func GenerateSpeechMP3(speed float64, voice, text, fileName string) error {
	switch config.Params.TTSAPI {
	case "openai":
		return sound.GenerateOpenaiSpeechMP3(config.Params.OpenaiAPIURL, config.Params.ApiKey, speed, voice, text, fileName)
	case "google":
		return sound.GenerateGoogleSpeechMP3(speed, voice, text, fileName)
	case "elevenlabs":
		return sound.GenerateElevenlabsSpeechMP3(config.Params.ElevenlabsAPIKey, voice, text, fileName)
	default:
		return fmt.Errorf("unsupported engine: %s", config.Params.TTSAPI)
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
