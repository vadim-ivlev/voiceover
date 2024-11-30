package epubs

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
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
func processContent(content, filePath string, epubTextLines []EpubTextLine) string {
	// Implement the actual content processing logic here
	// For now, just return the content unchanged
	return content
}
