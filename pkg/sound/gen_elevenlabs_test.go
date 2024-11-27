package sound

import (
	"testing"

	"github.com/vadim-ivlev/voiceover/internal/config"
)

func TestGenerateElevenlabsSpeechMP3(t *testing.T) {
	type args struct {
		apiKey   string
		voice    string
		text     string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test GenerateElevenlabsSpeechMP3 Sam",
			args: args{
				apiKey:   config.Params.ElevenlabsAPIKey,
				voice:    ElevenlabsVoiceSam,
				text:     "Hello, Sam!",
				fileName: "hello_sam.mp3",
			},
			wantErr: false,
		},
		{
			name: "Test GenerateElevenlabsSpeechMP3 Alice",
			args: args{
				apiKey:   config.Params.ElevenlabsAPIKey,
				voice:    ElevenlabsVoiceAlice,
				text:     "Hello, Alice!",
				fileName: "hello_alice.mp3",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateElevenlabsSpeechMP3(tt.args.apiKey, tt.args.voice, tt.args.text, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("GenerateElevenlabsSpeechMP3() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
