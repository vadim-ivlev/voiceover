package epubs

import (
	"archive/zip"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

// ListEpubFiles returns the content of the EPUB file treating it as a zip file
func ListEpubFiles(epubPath string) ([]string, error) {
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return nil, err
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
func filterExtension(files []string, ext string) []string {
	var filtered []string
	for _, f := range files {
		if filepath.Ext(f) == ext {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// listTranslatableFiles returns a slice of strings containing only the files that can be translated
func listTranslatableFiles(files []string) []string {
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
			defer rc.Close()

			return readAll(rc)
		}
	}
	return "", nil
}

// readAll reads all the content from the reader and returns it as a string
func readAll(r io.Reader) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
