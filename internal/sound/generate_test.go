package sound

import (
	"testing"
)

func TestGenerateMP3(t *testing.T) {
	type args struct {
		speed    float64
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
			name: "Test GenerateMP3",
			args: args{
				speed:    0.7,
				voice:    VoiceNova,
				text:     "The quick brown fox.",
				fileName: "The_quick_brown_fox.mp3",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateSpeechMP3(tt.args.speed, tt.args.voice, tt.args.text, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("GenerateMP3() error = %v, wantErr %v", err, tt.wantErr)
			}
			err := PlayMP3(tt.args.fileName)
			t.Errorf("PlayMP3() error = %v", err)
		})
	}
}

func TestGenerateSilenceFileFfmpeg(t *testing.T) {
	type args struct {
		fileName string
		duration float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test GenerateSilenceFileFfmpeg",
			args: args{
				fileName: "silence.mp3",
				duration: 0.5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateSilenceMP3(tt.args.fileName, tt.args.duration); (err != nil) != tt.wantErr {
				t.Errorf("GenerateSilenceFileFfmpeg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
