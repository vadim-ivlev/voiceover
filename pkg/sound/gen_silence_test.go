package sound

import "testing"

func TestGenerateSilenceMP3(t *testing.T) {
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
			if err := GenerateSilenceMP3(tt.args.duration, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("GenerateSilenceFileFfmpeg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
