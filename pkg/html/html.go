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

type HtmlTextLine struct {
	// The text content of the line
	Text string
	// HTML text of the line
	Html string
	// The line number in the file
	Index int
	// Selector used to extract the text line
	Selector string
}

// FetchSelectorTextsFromHTML extracts lines of text from HTML content.
// Each line of text is the text content of an HTML element selected by the provided CSS selector.
// The text lines are returned in the order of appearance in the HTML content.
// Parameters:
// - htmlContent: the HTML content to extract text lines from.
// - cssSelector: the CSS selector to select the HTML elements to extract text from.
// For example, "p", "h1", "h2", "h3", "div", "span", etc.
// Returns:
// - textLines: the extracted HTML text lines.
// - err: an error if the extraction failed.
func FetchSelectorTextsFromHTML(htmlContent, cssSelector string) (textLines []HtmlTextLine, err error) {
	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// Find elements matching the CSS selector and extract their text content
	doc.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			fmt.Printf("Failed to get HTML content for element %d:  %v\n", i, err)
		}

		htmlTextLine := HtmlTextLine{
			Text:     s.Text(),
			Html:     html,
			Index:    i,
			Selector: cssSelector,
		}
		textLines = append(textLines, htmlTextLine)
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

// UpdateHTMLWithSelectorEpubTextLines updates the text content of HTML elements
// matching a given CSS selector with the provided text lines.
//
// Parameters:
//   - htmlContent: The original HTML content as a string.
//   - cssSelector: The CSS selector to find the elements to be updated.
//   - textLines: A slice of HtmlTextLine structs containing the index and text
//     to update the elements with.
//
// Returns:
//   - modifiedHTML: The modified HTML content as a string.
//   - err: An error if any occurred during the parsing or updating process.
func UpdateHTMLWithSelectorEpubTextLines(htmlContent, cssSelector string, htmlTextLines []HtmlTextLine) (modifiedHTML string, err error) {
	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	// Find elements matching the CSS selector
	elements := doc.Find(cssSelector)

	// and replace their text content with the provided text lines
	for _, line := range htmlTextLines {

		selector := line.Selector
		// If the selector of the text line does not match the CSS selector,
		if selector != cssSelector {
			continue
		}

		text := strings.TrimSpace(line.Text)
		// if text is empty, skip the update
		if text == "" {
			continue
		}

		idx := line.Index
		// if the index is out of the range of elements, skip the update
		if idx >= elements.Length() || idx < 0 {
			continue
		}

		html := line.Html

		// find the element at the specified index
		selectedElement := elements.Eq(idx)
		// set the text content of the element
		// selectedElement.SetText(text)
		selectedElement.SetHtml(html)
	}

	modifiedHTML, err = doc.Html()

	return modifiedHTML, err
}

// GetTextFromHTML extracts and returns the text content from the provided HTML string.
//
// Parameters:
//   - htmlContent: A string containing the HTML content to be parsed.
//
// Returns:
//   - text: A string containing the extracted text content from the HTML.
//   - err: An error if there is an issue parsing the HTML content.
func GetTextFromHTML(htmlContent string) (text string, err error) {
	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	// Get the text content of the HTML document
	text = doc.Text()
	return text, nil
}
