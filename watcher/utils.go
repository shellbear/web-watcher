package watcher

import (
	"errors"

	"golang.org/x/net/html"
)

// A simple helper function to iterate over HTML node.
func crawlDocument(node *html.Node, handler func(node *html.Node) bool) bool {
	if handler(node) {
		return true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if crawlDocument(child, handler) {
			return true
		}
	}

	return false
}

// Extract all tags from HTML page.
func extractTags(doc *html.Node) []string {
	var tags []string
	crawlDocument(doc, func(node *html.Node) bool {
		if node.Type == html.ElementNode {
			tags = append(tags, node.Data)
		}

		return false
	})

	return tags
}

// Extract body from HTML page.
func getBody(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	crawlDocument(doc, func(node *html.Node) bool {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return true
		}

		return false
	})

	if body != nil {
		return body, nil
	}
	return nil, errors.New("missing <body> in the node tree")
}
