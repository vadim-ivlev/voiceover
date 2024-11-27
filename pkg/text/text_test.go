package text

import (
	"fmt"
	"os"
	"testing"

	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/pkg/logger"
)

func TestMain(m *testing.M) {
	// No colors
	logger.NoColor = true

	// change directory to the root of the project
	os.Chdir("../..")
	app.InitLoggerSetParams()
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

func TestSplitTextFileLines(t *testing.T) {
	type args struct {
		fileName   string
		startIndex int
		endIndex   int
	}
	tests := []struct {
		name          string
		args          args
		wantNumLines  int
		wantFirstLine string
		wantLastLine  string
		wantStart     int
		wantEnd       int
		wantErr       bool
	}{
		{
			name: "Test SplitTextFileLines Normal",
			args: args{
				fileName:   "./texts/dahl.txt",
				startIndex: 0,
				endIndex:   19,
			},
			wantNumLines:  19,
			wantFirstLine: "ROALD DAHL",
			wantLastLine:  "Contents",
			wantStart:     0,
			wantEnd:       19,
			wantErr:       false,
		},
		{
			name: "Test SplitTextFileLines negative StartIndex",
			args: args{
				fileName:   "./texts/dahl.txt",
				startIndex: -1,
				endIndex:   19,
			},
			wantNumLines:  19,
			wantFirstLine: "ROALD DAHL",
			wantLastLine:  "Contents",
			wantStart:     0,
			wantEnd:       19,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLines, gotStart, gotEnd, err := GetTextFileLines(tt.args.fileName, tt.args.startIndex, tt.args.endIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitTextFileLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStart != tt.wantStart {
				t.Errorf("SplitTextFileLines() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("SplitTextFileLines() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
			if len(gotLines) != tt.wantNumLines {
				t.Errorf("SplitTextFileLines() len(gotLines) = %v, want %v", len(gotLines), tt.wantNumLines)
			}
			if len(gotLines) > 0 {
				if gotLines[0] != tt.wantFirstLine {
					t.Errorf("SplitTextFileLines() gotLines[0] = %v, want %v", gotLines[0], tt.wantFirstLine)
				}
				if gotLines[len(gotLines)-1] != tt.wantLastLine {
					t.Errorf("SplitTextFileLines() gotLines[last] = %v, want %v", gotLines[len(gotLines)-1], tt.wantLastLine)
				}
			}
		})
	}
}
