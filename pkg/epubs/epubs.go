package epubs

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vadim-ivlev/voiceover/pkg/html"
	"github.com/vadim-ivlev/voiceover/pkg/utils"
)

// ListEpubFiles returns the content of the EPUB file treating it as a zip file
//
// Params:
//
//	epubPath - the path to the EPUB file
//
// Return:
//
//	a sorted list of file names within the EPUB and an error if any
func ListEpubFiles(epubPath string) ([]string, error) {
	// print current directory
	dir, err := os.Getwd()
	if err == nil {
		fmt.Printf("Current directory: %s\n", dir)
	}
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open EPUB file %s: %w", epubPath, err)
	}
	defer r.Close()

	var files []string
	for _, f := range r.File {
		files = append(files, f.Name)
	}
	sort.Strings(files)
	return files, nil
}

// filterExtension returns a slice of strings containing only the files with the specified extension
//
// Params:
//
//   - files - a list of file names
//   - ext - the extension to filter by
//
// Return:
//
//   - a filtered list of file names with the specified extension
func filterExtension(files []string, ext string) []string {
	var filtered []string
	for _, f := range files {
		if strings.EqualFold(filepath.Ext(f), ext) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// # listProcessableFiles
//
// returns a slice of strings containing only the files that can be translated
//
// Params:
//
//   - files - a list of file names
//
// Return:
//
//   - a list of file names with extensions .ncx, .xhtml, and .html
func listProcessableFiles(files []string) []string {
	translatableFiles := []string{}

	ncxs := filterExtension(files, ".ncx")
	sort.Strings(ncxs)
	translatableFiles = append(translatableFiles, ncxs...)

	xhtmls := filterExtension(files, ".xhtml")
	sort.Strings(xhtmls)
	translatableFiles = append(translatableFiles, xhtmls...)

	htmls := filterExtension(files, ".html")
	sort.Strings(htmls)
	translatableFiles = append(translatableFiles, htmls...)

	return translatableFiles
}

// getFileContent returns the content of the file with the specified name from the EPUB file
//
// Params:
//
//   - epubPath - the path to the EPUB file
//   - fileName - the name of the file to retrieve content from
//
// Return:
//
//   - the content of the file as a string and an error if any
func getFileContent(epubPath string, fileName string) (string, error) {
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == fileName {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			content, err := readAll(rc)
			rc.Close() // Close immediately after reading
			if err != nil {
				return "", err
			}
			return content, nil
		}
	}

	return "", nil
}

// readAll reads all the content from the reader and returns it as a string
//
// Params:
//
//   - r - an io.Reader to read from
//
// Return:
//
//   - the content read as a string and an error if any
func readAll(r io.Reader) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetAllEpubTextLines extracts all translatable text lines from an EPUB file.
//
// Parameters:
//   - epubPath: The file path to the EPUB file.
//   - selectors: A slice of CSS selectors used to identify the HTML elements to process.
//
// Returns:
//   - epubTexts: A slice of EpubTextLine containing the translatable text lines from the EPUB file.
//   - err: An error if any occurred during the process of extracting text lines.
//
// The function performs the following steps:
//  1. Lists all files in the EPUB archive.
//  2. Filters the files to identify those that contain translatable text.
//  3. Reads the content of each translatable file.
//  4. Extracts the translatable text lines from the content.
//  5. Returns a slice of EpubTextLine containing all the extracted text lines.
func GetAllEpubTextLines(epubPath string, selectors []string) (epubTexts []EpubTextLine, err error) {
	files, err := ListEpubFiles(epubPath)
	if err != nil {
		return
	}

	translatableFiles := listProcessableFiles(files)
	epubTextLines := []EpubTextLine{}
	for _, f := range translatableFiles {

		content, err := getFileContent(epubPath, f)
		if err != nil {
			return nil, err
		}

		translatableEpubTextLines, err := fetchProcessableLines(f, content, selectors)
		if err != nil {
			return nil, err
		}
		epubTextLines = append(epubTextLines, translatableEpubTextLines...)

	}

	return epubTextLines, nil
}

// GetEpubTextLines extracts a range of text lines from an EPUB file based on the provided start and end indexes and CSS selectors.
//
// Parameters:
//   - epubPath: The file path to the EPUB file.
//   - startIndex: The starting index of the text lines to extract.
//   - endIndex: The ending index of the text lines to extract.
//   - selectors: A slice of CSS selectors to filter the text lines.
//
// Returns:
//   - epubTexts: A slice of EpubTextLine containing the extracted text lines.
//   - start: The actual starting index used for extraction.
//   - end: The actual ending index used for extraction.
//   - err: An error object if an error occurred during the extraction process, otherwise nil.
func GetEpubTextLines(epubPath string, startIndex, endIndex int, selectors []string) (epubTexts []EpubTextLine, start, end int, err error) {
	allEpubTexts, err := GetAllEpubTextLines(epubPath, selectors)
	if err != nil {
		return nil, 0, 0, err
	}

	// calculate the start and end indexes
	start, end, err = utils.CalcStartEndIndex(len(allEpubTexts), startIndex, endIndex)
	if err != nil {
		return nil, 0, 0, err
	}

	// get the lines
	epubTexts = allEpubTexts[start:end]
	return epubTexts, start, end, nil
}

// fetchSelectorLines returns a slice of EpubTextLine objects representing the text content of the HTML elements
// selected by the specified CSS selector
//
// Params:
//
//   - epubPath - the path to the file in the EPUB
//   - content - the content of the file
//   - cssSelector - the CSS selector to select the HTML elements
//
// Return:
//
//   - a slice of EpubTextLine objects and an error if any
func fetchSelectorLines(epubPath, content, cssSelector string) (epubTextLines []EpubTextLine, err error) {
	texts, err := html.FetchSelectorTextsFromHTML(content, cssSelector)
	if err != nil {
		return nil, err
	}
	epubTextLines = []EpubTextLine{}
	for i, text := range texts {
		epubTextLines = append(epubTextLines, EpubTextLine{
			Text:     text,
			Index:    i,
			FilePath: epubPath,
			Selector: cssSelector,
		})
	}
	return epubTextLines, nil
}

// fetchProcessableLines returns a slice of EpubTextLine objects representing the text content of the file with the specified name
//
// Params:
//
//   - epubPath - the path to the file in the EPUB
//   - content - the content of the file
//   - selectors - the list of CSS selectors to select the HTML elements
//
// Return:
//
//   - a slice of EpubTextLine objects and an error if any
func fetchProcessableLines(epubPath, content string, selectors []string) (epubTextLines []EpubTextLine, err error) {
	epubTextLines = []EpubTextLine{}

	// loop trough selectors to fetch text lines
	for _, selector := range selectors {
		lines, err := fetchSelectorLines(epubPath, content, selector)
		if err != nil {
			return nil, err
		}
		epubTextLines = append(epubTextLines, lines...)
	}

	return epubTextLines, nil
}
