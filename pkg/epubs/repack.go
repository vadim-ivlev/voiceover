package epubs

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vadim-ivlev/voiceover/pkg/html"
)

// RepackEpub repacks an existing EPUB file by copying its contents to a new EPUB file,
// while allowing modifications to the EPUB content using the provided EpubTextLine structs.
//
// Parameters:
//   - existingEpubPath: The file path to the existing EPUB file.
//   - updatedEpubPath: The file path where the new EPUB file will be created.
//   - epubTextLines: A slice of EpubTextLine used to rewrite the content of the file if not empty.
//
// Returns:
//   - error: An error if any occurs during the repacking process.
//
// The function performs the following steps:
//  1. Opens the existing EPUB file as a ZIP archive.
//  2. Creates a new ZIP file for the updated EPUB.
//  3. Copies all files from the old ZIP archive to the new one, applying modifications if necessary.
//  4. Logs and returns any errors encountered during the process.
func RepackEpub(existingEpubPath, updatedEpubPath string, epubTextLines []EpubTextLine) error {

	// Open the existing ZIP archive
	oldZip, err := zip.OpenReader(existingEpubPath)
	if err != nil {
		log.Fatalf("Failed to open existing ZIP file: %v", err)
	}
	defer oldZip.Close()

	// Create a new ZIP file
	newZip, err := os.Create(updatedEpubPath)
	if err != nil {
		log.Fatalf("Failed to create new ZIP file: %v", err)
	}
	defer newZip.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(newZip)
	defer zipWriter.Close()

	// Copy all files from the old ZIP archive to the new one
	for _, file := range oldZip.File {
		err := addFileToArchive(file, zipWriter, epubTextLines)
		if err != nil {
			log.Fatalf("Failed to copy file '%s': %v", file.Name, err)
		}
	}

	fmt.Printf("Successfully repacked ZIP archive '%s' to  '%s'.\n", existingEpubPath, updatedEpubPath)
	return nil
}

// addFileToArchive copies a file from a source ZIP archive to a destination ZIP archive.
// If epubTextLines is not empty, it rewrites the content of the file using the provided lines.
//
// Parameters:
//   - file: The source file to be copied.
//   - writer: The destination ZIP writer where the file will be copied to.
//   - epubTextLines: A slice of EpubTextLine used to rewrite the content of the file if not empty.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func addFileToArchive(file *zip.File, writer *zip.Writer, epubTextLines []EpubTextLine) error {
	srcFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %v", file.Name, err)
	}
	defer srcFile.Close()

	// Create a new file in the destination ZIP archive
	destFileWriter, err := writer.Create(file.Name)
	if err != nil {
		return fmt.Errorf("failed to create file '%s' in new ZIP: %v", file.Name, err)
	}

	// Rewrite the content of the file if epubTextLines is not empty
	if len(epubTextLines) > 0 {
		_, err = processAndWrite(destFileWriter, srcFile, file.Name, epubTextLines)
		if err != nil {
			return fmt.Errorf("failed to rewrite content of file '%s': %v", file.Name, err)
		}
	} else {
		// If no epubTextLines are provided, copy the content as is
		_, err = io.Copy(destFileWriter, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy content of file '%s': %v", file.Name, err)
		}
	}

	return nil
}

// processAndWrite reads from the provided src Reader, processes the content based on the given
// epubTextLines, and writes the result to the dest Writer. It returns the number of bytes
// written and any error encountered during the process.
//
// Parameters:
//
//	dest - the destination Writer where the processed content will be written
//	src - the source Reader from which the content will be read
//	epubPath - the path of the file in the EPUB being processed
//	epubTextLines - a slice of EpubTextLine that contains the text lines to be processed
//
// Returns:
//
//	The number of bytes written to the dest Writer and an error if any occurred during the process.
func processAndWrite(dest io.Writer, src io.Reader, epubPath string, epubTextLines []EpubTextLine) (int64, error) {
	// Read the entire content from the source
	srcContent, err := io.ReadAll(src)
	if err != nil {
		return 0, fmt.Errorf("failed to read source content: %v", err)
	}

	// Process the content
	processedContent := processContent(string(srcContent), epubPath, epubTextLines)

	// Write the processed content to the destination
	bytesWritten, err := io.WriteString(dest, processedContent)
	if err != nil {
		return 0, fmt.Errorf("failed to write processed content: %v", err)
	}

	return int64(bytesWritten), nil
}

// processContent processes the given content based on the provided file path and EPUB text lines.
// It returns the processed content as a string.
//
// Parameters:
//   - content: The original content to be processed.
//   - filePath: The path to the file associated with the content.
//   - epubTextLines: A slice of EpubTextLine structs representing the lines of text in the EPUB.
//
// Returns:
//   - A string containing the processed content.
func processContent(content, filePath string, epubTextLines []EpubTextLine) (processedContent string) {
	// Filter lines for the current file
	relevantLines := filterLinesByFilePath(epubTextLines, filePath)

	// for each prcessable selector, modify the content
	for _, selector := range ProcessableSelectors {
		// Filter lines by selector
		selectorLines := filerLinesBySelector(relevantLines, selector)
		if len(selectorLines) > 0 {
			//Update the content with the selector lines
			content = processContentBySelector(content, selector, selectorLines)
		}
	}

	return content
}

// filterLinesByFilePath filters the given slice of EpubTextLine and returns a new slice
// containing only the lines that match the specified file path.
//
// Parameters:
//   - epubTextLines: A slice of EpubTextLine structs to be filtered.
//   - filePath: The file path to filter the lines by.
//
// Returns:
//
//	A slice of EpubTextLine structs that have the specified file path.
func filterLinesByFilePath(epubTextLines []EpubTextLine, filePath string) []EpubTextLine {
	var filteredLines []EpubTextLine
	for _, line := range epubTextLines {
		if line.FilePath == filePath {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
}

// filerLinesBySelector filters the given slice of EpubTextLine and returns a new slice
// containing only the lines that match the specified selector.
//
// Parameters:
//   - epubTextLines: A slice of EpubTextLine structs to be filtered.
//   - selector: The selector to filter the lines by.
//
// Returns:
//
//	A slice of EpubTextLine structs that have the specified selector.
func filerLinesBySelector(epubTextLines []EpubTextLine, selector string) []EpubTextLine {
	var filteredLines []EpubTextLine
	for _, line := range epubTextLines {
		if line.Selector == selector {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
}

// processContentBySelector processes the given HTML content by updating elements
// that match the specified CSS selector with the provided text lines.
//
// Parameters:
//   - htmlContent: The original HTML content to be modified.
//   - cssSelector: The CSS selector used to identify elements in the HTML content.
//   - epubTextLines: A slice of EpubTextLine structs containing the text lines to be inserted.
//
// Returns:
//   - modifiedHTML: The modified HTML content with the text lines inserted. If an error occurs
//     during the update, the original HTML content is returned.
//
// This function uses the UpdateHTMLWithSelectorTexts function from the html package to update the HTML content.
// Before calling the function, it extracts the text lines from the EpubTextLine structs and passes them to the function.
func processContentBySelector(htmlContent, cssSelector string, epubTextLines []EpubTextLine) (modifiedHTML string) {
	// Convert the HtmlTextLine slice from the EpubTextLine slice
	htmlTextLines := epubToHtmlTextLines(epubTextLines)
	modifiedHTML, err := html.UpdateHTMLWithSelectorEpubTextLines(htmlContent, cssSelector, htmlTextLines)
	if err != nil {
		return htmlContent
	}
	return modifiedHTML
}

// epubToHtmlTextLines converts a slice of EpubTextLine to a slice of html.HtmlTextLine.
// Each EpubTextLine is mapped to an html.HtmlTextLine with the same Text, Index, and Selector fields.
//
// Parameters:
//   - epubTextLines: A slice of EpubTextLine structs to be converted.
//
// Returns:
//   - A slice of html.HtmlTextLine structs with the corresponding fields from the input slice.
func epubToHtmlTextLines(epubTextLines []EpubTextLine) []html.HtmlTextLine {
	htmlTextLines := make([]html.HtmlTextLine, len(epubTextLines))
	for i, epubLine := range epubTextLines {
		htmlTextLines[i] = html.HtmlTextLine{
			Text:     epubLine.Text,
			Index:    epubLine.Index,
			Selector: epubLine.Selector,
		}
	}
	return htmlTextLines
}
