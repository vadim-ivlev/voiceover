// Translates table of contents files inside the EPUB file

package epubs

import (
	"fmt"
	"path"
	"time"

	"github.com/vadim-ivlev/voiceover/internal/config"
	"github.com/vadim-ivlev/voiceover/pkg/translator"
)

// TablesOfContents is a map of file paths of .ncx files inside epub to their translated contents.
var TablesOfContents = map[string]string{}

// TranslateTablesOfContents translates the table of contents (TOC) of an EPUB file.
// It takes the path to the EPUB file as input and returns a map where the keys are the TOC file names
// and the values are the translated contents of those TOC files. If the provided file is not an EPUB,
// or if there is an error during processing, it returns an error.
//
// Parameters:
//   - epubPath: The path to the EPUB file.
//
// Returns:
//   - tocs: A map containing the translated TOC contents, where the keys are the TOC file names
//     and the values are the translated contents.
//   - err: An error if any occurred during the translation process.
func TranslateTablesOfContents(epubPath string) (tocs map[string]string, err error) {
	tocs = make(map[string]string)

	// Check if the file extension is .epub
	if path.Ext(epubPath) != ".epub" {
		return nil, fmt.Errorf("file is not an EPUB")
	}

	// List all files in the EPUB
	files, err := ListEpubFiles(epubPath)
	if err != nil {
		return nil, fmt.Errorf("error listing EPUB files: %w", err)
	}

	// Filter files to process only .ncx files
	processableFiles := listProcessableFiles(files, []string{".ncx"})

	for _, file := range processableFiles {
		// Read the content of the file
		contents, err := getFileContent(epubPath, file)
		if err != nil {
			fmt.Printf("Error reading file content: %v\n", err)
			continue
		}

		// Log the start of translation
		fmt.Printf("--------------- Translating TOC %s started\n", file)
		startTime := time.Now()

		// Translate the content
		translatedContents, err := translator.TranslateText(config.Params.OpenaiAPIURL, config.Params.ApiKey, config.Params.TranslateTo, translator.TocTranslInstructions, contents)
		if err != nil {
			return nil, fmt.Errorf("error translating TOC %s: %w", file, err)
		}

		// Log the time taken for translation
		fmt.Printf("--------------- Translation of TOC %s took %v\n", file, time.Since(startTime))

		// Store the translated content
		tocs[file] = translatedContents
	}

	// Update the global map with the translated TOCs
	TablesOfContents = tocs
	return tocs, nil
}
