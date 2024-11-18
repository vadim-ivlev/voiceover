package pipe

import (
	"reflect"
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

func TestStartPipeline(t *testing.T) {
	type args struct {
		textLines []string
	}
	tests := []struct {
		name         string
		args         args
		wantDoneJobs []Job
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDoneJobs, err := StartPipeline(tt.args.textLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartPipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDoneJobs, tt.wantDoneJobs) {
				t.Errorf("StartPipeline() = %v, want %v", gotDoneJobs, tt.wantDoneJobs)
			}
		})
	}
}
