package texts

import (
	"bufio"
	"bytes" // added
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/voiceover/pkg/utils"
)

var splitString = "\n"

// splitTextFile splits a text file into an array of strings.
func splitTextFile(fileName string) ([]string, error) {
	// read the text file
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	text := string(bytes)
	// split the text file with separator "\n\n"
	lines := strings.Split(text, splitString)

	return lines, nil
}

// splitTextFileScan splits a text file into an array of strings.
func splitTextFileScan(fileName string) (lines []string, err error) {
	// read the text file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(split)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// split is a split function for a Scanner that returns each line of text, stripped of any trailing end-of-line marker.
func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte(splitString)); i >= 0 {
		return i + len(splitString), data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

// SaveTextFile saves a text to a file.
func SaveTextFile(fileName, text string) error {
	// write the text to the file
	err := os.WriteFile(fileName, []byte(text), 0644)
	if err != nil {
		log.Error().Msgf("Failed to save text file %s: %v", fileName, err)
		return err
	}
	return nil
}

// SaveTextFileLines saves an array of strings to a text file.
func SaveTextFileLines(fileName string, lines []string) error {
	return SaveTextFile(fileName, strings.Join(lines, "\n"))
}

// GetTextFileLines - splits a text file into an array of strings.
// Parameters:
// - fileName: the name of the file to split.
// - startIndex: the index of the first line to include.
// - endIndex: the index of the last line. The last line will not be included.
// Returns:
// - lines: the array of strings. Can be empty.
// - start: the actual index of the first line of the returned strings.
// - end: the actual index of the last line of the returned strings.
func GetTextFileLines(fileName string, startIndex, endIndex int) (lines []string, start, end int, err error) {
	allLines, err := splitTextFileScan(fileName)
	if err != nil {
		return nil, 0, 0, err
	}

	// calculate the start and end indexes
	start, end, err = utils.CalcStartEndIndex(len(allLines), startIndex, endIndex)
	if err != nil {
		return nil, 0, 0, err
	}

	// get the lines
	lines = allLines[start:end]
	return lines, start, end, nil
}
