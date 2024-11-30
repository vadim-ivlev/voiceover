package epubs

import (
	"fmt"
	"log"

	// "github.com/bmaupin/go-epub"
	"github.com/medialibraryonline/epub-go"
)

func main() {
	// Open the EPUB
	epubPath := "texts/dahl.epub"
	book, err := epub.Open(epubPath)
	if err != nil {
		log.Fatalf("Failed to open EPUB: %v", err)
	}

	// Get the title
	title := book.Opf.Spine
	fmt.Printf("Title: %s\n", title)

	// Get the content
	chapters := book.Files()
	for i, chapter := range chapters {
		fmt.Printf("#%d %s\n", i+1, chapter)
	}
}
