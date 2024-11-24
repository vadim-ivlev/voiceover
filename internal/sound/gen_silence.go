package sound

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/xfrr/goffmpeg/transcoder"
)

// GenerateSilenceMP3 - creates a  mp3 file that contains silence using ffmpeg.
// Parameters:
// fileName - the name of the file to create.
// duration - the duration of the silence in seconds.
func GenerateSilenceMP3(duration float64, fileName string) error {
	t := fmt.Sprintf("%f", duration)
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "anullsrc=r=24000:cl=mono", "-t", t, "-q:a", "9", "-acodec", "libmp3lame", fileName)
	err := cmd.Run()
	return err
}

func GenerateSilenceMP31(filename string, durationSeconds int) error {
	// Define sample rate and total samples for silence
	sampleRate := 44100
	totalSamples := sampleRate * durationSeconds

	// Create a buffer of zeroes (silence)
	silence := make([]byte, totalSamples*2) // 16-bit audio

	// Write silence to a WAV file first
	wavFile, err := os.Create("temp.wav")
	if err != nil {
		return err
	}
	defer wavFile.Close()

	// Write WAV header
	wavHeader := []byte{
		'R', 'I', 'F', 'F', // ChunkID
		0, 0, 0, 0, // ChunkSize (to be filled later)
		'W', 'A', 'V', 'E', // Format
		'f', 'm', 't', ' ', // Subchunk1ID
		16, 0, 0, 0, // Subchunk1Size (16 for PCM)
		1, 0, // AudioFormat (1 for PCM)
		1, 0, // NumChannels (1 for mono)
		byte(sampleRate & 0xff), byte((sampleRate >> 8) & 0xff), byte((sampleRate >> 16) & 0xff), byte((sampleRate >> 24) & 0xff), // SampleRate
		byte((sampleRate * 2) & 0xff), byte(((sampleRate * 2) >> 8) & 0xff), byte(((sampleRate * 2) >> 16) & 0xff), byte(((sampleRate * 2) >> 24) & 0xff), // ByteRate
		2, 0, // BlockAlign
		16, 0, // BitsPerSample
		'd', 'a', 't', 'a', // Subchunk2ID
		byte((totalSamples * 2) & 0xff), byte(((totalSamples * 2) >> 8) & 0xff), byte(((totalSamples * 2) >> 16) & 0xff), byte(((totalSamples * 2) >> 24) & 0xff), // Subchunk2Size
	}

	_, err = wavFile.Write(wavHeader)
	if err != nil {
		return err
	}

	// Write silence data
	_, err = wavFile.Write(silence)
	if err != nil {
		return err
	}

	// Convert WAV to MP3 using goffmpeg
	trans := new(transcoder.Transcoder)
	err = trans.Initialize("temp.wav", filename)
	if err != nil {
		return err
	}

	done := trans.Run(true)
	if err := <-done; err != nil {
		return err
	}

	return nil
}
