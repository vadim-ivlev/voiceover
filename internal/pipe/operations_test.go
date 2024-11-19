package pipe

import (
	"testing"
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
			args:    args{filePath: "./texts/The Mind.txt"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProcessFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("StartFileProcessing() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
