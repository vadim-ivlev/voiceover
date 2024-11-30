// Provide functions to extract text lines from HTML,
// and to modify the HTML using a provided list of text lines.
// The functions should use the "github.com/PuerkitoBio/goquery" package to parse the HTML content.
// The general usage flow is as follows:
// 1. Extract text lines from the HTML content using the provided CSS selector.
// 2. Modify the text lines as needed.
// 3. Modify the HTML content by replacing the text content of the HTML elements selected by the CSS selector with the modified text lines.

package html

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FetchSelectorTextsFromHTML extracts lines of text from HTML content.
// Each line of text is the text content of an HTML element selected by the provided CSS selector.
// The text lines are returned in the order of appearance in the HTML content.
// Parameters:
// - htmlContent: the HTML content to extract text lines from.
// - cssSelector: the CSS selector to select the HTML elements to extract text from.
// For example, "p", "h1", "h2", "h3", "div", "span", etc.
// Returns:
// - textLines: the extracted text lines.
// - err: an error if the extraction failed.
func FetchSelectorTextsFromHTML(htmlContent, cssSelector string) (textLines []string, err error) {
	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	// Find elements matching the CSS selector and extract their text content
	doc.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		textLines = append(textLines, s.Text())
	})
	return textLines, nil
}

// UpdateHTMLWithSelectorTexts modifies the HTML content by replacing the text content of the HTML
// elements selected by the provided CSS selector with the provided text lines.
// The text lines are used in the order of appearance in the HTML content.
// Parameters:
// - htmlContent: the HTML content to modify.
// - cssSelector: the CSS selector to select the HTML elements to modify.
// - textLines: the text lines to replace the text content of the selected HTML elements.
// Returns:
// - modifiedHTML: the modified HTML content.
// - err: an error if the modification failed.
func UpdateHTMLWithSelectorTexts(htmlContent, cssSelector string, textLines []string) (modifiedHTML string, err error) {
	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}
	index := 0
	// Find elements matching the CSS selector and replace their text content with the provided text lines
	doc.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		if index < len(textLines) {
			s.SetText(textLines[index])
			index++
		}
	})

	// Serialize the modified HTML content
	modifiedHTML, err = doc.Html()
	return modifiedHTML, err
}

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
	fmt.Printf("Extracted text lines: \n%#v\n", textLines)

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
