package epubs

// EpubTextLine represents a line of text extracted from an EPUB file.
// It contains the text content, the line number, the file path, the selector used for extraction,
// and a flag indicating whether the node has a child node.
type EpubTextLine struct {
	// The text content of the line
	Text string
	// The line number in the file
	Index int
	// Path to the file in the EPUB
	FilePath string
	// Selector used to extract the text line
	Selector string
	// // If the node has a child node
	// HasChildren bool
}

// ProcessableSelectors contains the list of selectors that can be processed.
var ProcessableSelectors = []string{"h1", "h2", "h3", "h4", "h5", "h6", "p", "div"}
