package epubs

import (
	"fmt"
	"testing"

	"github.com/vadim-ivlev/voiceover/pkg/utils"
)

var HTMLText = `

	<p>First 
	
	paragraph</p>
	<p>Second <a href="#">some link</a> paragraph</p>
	<div>
		what can <b>we</b> do?
	</div>

`

func Test_extractTextChunksXNet(t *testing.T) {
	type args struct {
		htmlText string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test extractTextChunks",
			args: args{
				htmlText: HTMLText,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChunks, err := extractTextChunksXNet(tt.args.htmlText)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTextChunks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(utils.PrettyJSON(gotChunks))
		})
	}
}

func Test_extractTextChunksGoquery(t *testing.T) {
	type args struct {
		htmlText string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test extractTextChunks",
			args: args{
				htmlText: HTMLText,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChunks, err := extractTextChunksGoquery(tt.args.htmlText)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTextChunksGoquery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(utils.PrettyJSON(gotChunks))
		})
	}
}
