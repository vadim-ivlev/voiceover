package html

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// TextChunk represents a chunk of text and its CSS selector.
type TextChunk struct {
	Index    int
	Text     string
	Selector string
}

// extractTextChunksXNet extracts text chunks from an HTML document using the x/net/html package.
func extractTextChunksXNet(htmlText string) (chunks []TextChunk, err error) {
	doc, err := html.Parse(strings.NewReader(htmlText))
	if err != nil {
		return nil, err
	}

	var textChunks []TextChunk
	var f func(*html.Node, string, int) int
	f = func(n *html.Node, selector string, index int) int {
		if n.Type == html.TextNode {
			textChunks = append(textChunks, TextChunk{Text: n.Data, Selector: selector, Index: index})
			index++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			childSelector := selector
			if c.Type == html.ElementNode {
				childSelector = selector + "." + c.Data
			}
			index = f(c, childSelector, index)
		}
		return index
	}
	f(doc, "", 0)

	return textChunks, nil
}

// extractTextChunksGoquery extracts text chunks from an HTML document using the goquery package.
func extractTextChunksGoquery(htmlText string) (chunks []TextChunk, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, err
	}

	var textChunks []TextChunk

	selection := doc.Find("a")

	selection.Each(func(i int, s *goquery.Selection) {

		// fmt.Printf("NodeName: %s\n", goquery.NodeName(s))
		// fmt.Printf("Text: %s\n", s.Text())
		// fmt.Println("=====================================")
		// fmt.Printf("Nodes[0] = %#v\n", s.Nodes[0])
		// if s.Nodes[0].Type == html.ElementNode {
		text := s.Text()
		textChunks = append(textChunks, TextChunk{
			Index:    i,
			Text:     text,
			Selector: goquery.NodeName(s),
		})
		s.SetText(text + " - modified")
		// }
	})
	fmt.Println(doc.Html())
	return textChunks, nil
}
