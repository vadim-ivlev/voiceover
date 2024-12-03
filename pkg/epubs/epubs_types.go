package epubs

// EpubTextLine represents a line of text extracted from an EPUB file.
// It contains the text content, the line number, the file path, the selector used for extraction,
// and a flag indicating whether the node has a child node.
type EpubTextLine struct {
	// The text content of the line
	Text string
	// HTML text of the line
	Html string
	// The index of the line in the HTML content
	htmlTextLineIndex int

	// The line number in the file
	Index int
	// Path to the file in the EPUB
	FilePath string
	// Selector used to extract the text line
	Selector string
	// // If the node has a child node
	// HasChildren bool
}

// ProcessableExtensions contains the list of file extensions that can be processed.
// ".ncx",
var ProcessableExtensions = []string{".html"}

// ProcessableSelectors contains the list of selectors that can be processed.
var ProcessableSelectors = []string{"blockquote>span", "h1", "h2", "h3", "h4", "p"}
