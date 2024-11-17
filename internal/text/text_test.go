package text

import (
	"fmt"
	"os"
	"testing"

	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/internal/logger"
)

func TestMain(m *testing.M) {
	// No colors
	logger.NoColor = true

	// change directory to the root of the project
	os.Chdir("../..")
	app.InitApp()
	fmt.Println("TestMain")

	os.Exit(m.Run())
}

func TestSplitTextFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Test SplitTextFile",
			args: args{
				fileName: "./texts/The Mind.txt",
			},
			want:    []string{"The quick brown fox.", "Jumps over the lazy dog."},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitTextFileScan(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitTextFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, v := range got {
				fmt.Println(i, v)
			}
		})
	}
}
