package pipe

import (
	"testing"

	"github.com/vadim-ivlev/voiceover/internal/sound"
)

func TestProcessFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test StartFileProcessing",
			args:    args{filePath: "./texts/The Mind is Flat ( PDFDrive ).txt"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProcessFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("StartFileProcessing() error = %v, wantErr %v", err, tt.wantErr)
			}
			sound.PlayMP3(tt.args.filePath + ".mp3")
		})
	}
}
