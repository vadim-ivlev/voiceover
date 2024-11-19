package main

import (
	"fmt"

	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/sound"
)

func main() {

	app.InitApp()

	speed := 0.7
	voice := sound.VoiceNova
	text := "The quick brown fox jumped over the lazy dog."
	fileName := "speech.mp3"

	err := sound.GenerateSpeechMP3(speed, voice, text, fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("MP3 file generated successfully.")
	}
}
