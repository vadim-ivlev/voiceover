package epubs

import (
	"fmt"
	"os"
	"testing"

	"github.com/vadim-ivlev/voiceover/internal/app"
	"github.com/vadim-ivlev/voiceover/pkg/logger"
	"github.com/vadim-ivlev/voiceover/pkg/utils"
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

func TestListEpubFiles(t *testing.T) {
	type args struct {
		epubPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			name: "ListEpubFiles",
			args: args{
				epubPath: "texts/dahl.epub",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListEpubFiles(tt.args.epubPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListEpubFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			translatableFiles := listProcessableFiles(got)
			fmt.Printf("translatableFiles = \n%s\n", utils.PrettyJSON(translatableFiles))

			ncxContent, err := getFileContent(tt.args.epubPath, translatableFiles[0])
			if err != nil {
				t.Errorf("getFileContent() error = %v", err)
				return
			}
			fmt.Printf("ncxContent = \n%s\n", ncxContent)

		})
	}
}

func TestGetEpubTextLines(t *testing.T) {
	type args struct {
		epubPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetEpubTextLines",
			args: args{
				epubPath: "texts/dahl.epub",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEpubTexts, err := GetEpubTextLines(tt.args.epubPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEpubTextLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("gotEpubTexts = \n%s\n", utils.PrettyJSON(gotEpubTexts[0:10]))
		})
	}
}
