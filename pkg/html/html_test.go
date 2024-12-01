// provide functions to extract text lines  from an HTML,
// and to modify the HTML using provided list of text lines.
// The functions should use the "github.com/PuerkitoBio/goquery" package to parse the HTML content.
// The general usage flow is as follows:
// 1. Extract text lines from the HTML content using the provided CSS selector.
// 2. Modify the text lines as needed.
// 3. Modify the HTML content by replacing the text content of the HTML elements selected by the CSS selector with the modified text lines.

package html

import (
	"fmt"
	"testing"

	"github.com/vadim-ivlev/voiceover/pkg/utils"
)

// Example usage
func usageExample() {
	htmlContent := `
        <html>
		<!-- This is a comment -->
            <body>
				<h1>Heading 1</h1>
                <h2>Heading 2</h2>
				<p></p>
				<p>  </p>
                <p>Paragraph 1 </p>
                <p>  Paragraph 2</p>
				<p><span>span</span></p>
            </body>
        </html>
    `
	cssSelector := "p"

	// Extract text lines from the HTML content
	textLines, err := FetchSelectorTextsFromHTML(htmlContent, cssSelector)
	if err != nil {
		fmt.Println("Failed to extract text lines from HTML content:", err)
		return
	}
	js := utils.PrettyJSON(textLines)
	fmt.Println("Extracted text lines: ")
	fmt.Println(js)

	// Define new text lines to replace the existing ones
	modifiedTextLines := []string{"", "  ", "NEW paragraph 1", "NEW paragraph 2", "NEW span", "NEW paragraph 3"}

	// Modify the HTML content with the new text lines
	modifiedHTML, err := UpdateHTMLWithSelectorTexts(htmlContent, cssSelector, modifiedTextLines)
	if err != nil {
		fmt.Println("Failed to modify HTML content with text lines:", err)
		return
	}
	fmt.Printf("Modified HTML content:\n%v\n", modifiedHTML)
}

func Test_usageExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test UsageExample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usageExample()
		})
	}
}
