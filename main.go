package main

import (
	"fmt"

	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/pkg/sound"
)

func main() {

	app.InitApp()

	text := "The quick brown fox jumped over the lazy dog."
	fileName := "speech.mp3"

	err := sound.GenerateMP3(text, fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("MP3 file generated successfully.")
	}
}
