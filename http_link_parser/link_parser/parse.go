package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// link in a html doc
type Link struct {
	Href string
	text string
}

func Parse(reader io.Reader) ([]Link, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)
	var links []Link

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildLink(node *html.Node) Link {
	var ret Link

	for _, attribute := range node.Attr {
		if attribute.Key == "href" {
			ret.Href = attribute.Val
			break
		}
	}

	ret.text = findText(node)

	return ret
}

func findText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	if node.Type != html.ElementNode {
		return ""
	}

	var return_value string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		return_value += findText(child)
	}

	return strings.Join(strings.Fields(return_value), "")
}

func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}

	var return_value []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		return_value = append(return_value, linkNodes(child)...)
	}

	return return_value
}
