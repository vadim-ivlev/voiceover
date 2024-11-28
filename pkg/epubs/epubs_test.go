package epubs

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

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test main",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
