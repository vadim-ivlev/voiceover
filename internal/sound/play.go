package sound

// import (
// 	"os"
// 	"time"

// 	"github.com/faiface/beep"
// 	"github.com/faiface/beep/mp3"
// 	"github.com/faiface/beep/speaker"
// )

// // PlayMP3 - plays an MP3 file.
// func PlayMP3(fileName string) error {
// 	// Open the MP3 file
// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Decode the MP3 file
// 	streamer, format, err := mp3.Decode(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer streamer.Close()

// 	// Initialize the speaker with the sample rate of the MP3 file
// 	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
// 	if err != nil {
// 		return err
// 	}

// 	// Play the MP3 file
// 	done := make(chan bool)
// 	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
// 		done <- true
// 	})))

// 	// Wait for the playback to finish
// 	<-done
// 	return nil
// }
