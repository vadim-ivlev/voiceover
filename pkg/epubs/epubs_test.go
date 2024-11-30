package epubs

import (
	"fmt"
	"testing"

	"github.com/vadim-ivlev/voiceover/pkg/utils"
)

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
			translatableFiles := listTranslatableFiles(got)
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
