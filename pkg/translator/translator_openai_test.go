package translator

import (
	"fmt"
	"testing"

	"github.com/vadim-ivlev/voiceover/internal/config"
)

func TestTranslateText(t *testing.T) {
	type args struct {
		apiURL   string
		apiKey   string
		language string
		text     string
	}
	tests := []struct {
		name            string
		args            args
		wantTranslation string
		wantErr         bool
	}{
		{
			name: "Test Hello World",
			args: args{
				apiURL:   config.Params.OpenaiAPIURL,
				apiKey:   config.Params.ApiKey,
				language: "German",
				text:     "Hello World",
			},
			wantTranslation: "Привет, мир",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTranslation, err := TranslateText(tt.args.apiURL, tt.args.apiKey, tt.args.language, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranslateText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotTranslation) == 0 {
				t.Errorf("TranslateText() gotTranslation length = 0, want length > 0")
			}
			fmt.Println(gotTranslation)
		})
	}
}
