package template

import (
	"bytes"
	"errors"
	"fmt"
	"go-blog/pkg/util/model"
	"io"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	// DateFormat is the layout used to parse from template
	DateFormat = "2006-01-02 15:04:05"
)

// ParseTemplate recives a string with HTML content and parse it to extract
// blog post metadata
func ParseTemplate(htmlContent string) (post model.Post, err error) {

	reader := strings.NewReader(htmlContent)
	doc, errParsing := html.Parse(reader)
	if errParsing != nil {
		err = fmt.Errorf("error parsing HTML content: %s", errParsing)
		return
	}

	// get head node
	head := extractHTMLNode(doc, "head")
	if head == nil {
		err = errors.New("HTML tag <head> was not found")
		return
	}

	// extract meta tags
	tags := extractMetaTags(head)
	if tags == nil {
		err = errors.New("no <meta> tags were found")
		return
	}

	for k, v := range tags {
		switch k {
		case "title":
			post.Title = v
		case "author":
			post.Author = v
		case "categories":
			post.Categories = v
		case "tags":
			post.Tags = v
		case "post-date":
			if parsedDate, errParse := time.Parse(DateFormat, v); errParse == nil {
				post.DateCreated = parsedDate
			}
		case "edit-date":
			if parsedDate, errParse := time.Parse(DateFormat, v); errParse == nil {
				post.DateUpdated = parsedDate
			}
		}
	}

	// get body node
	body := extractHTMLNode(doc, "body")
	if head == nil {
		err = errors.New("HTML tag <body> was not found")
		return
	}

	bodyString := nodeToString(body)
	post.Content = bodyString

	return
}

// extract meta tags values from head and put them into a map
func extractMetaTags(head *html.Node) (keys map[string]string) {
	for child := head.FirstChild; child != nil; child = child.NextSibling {

		if child.Data == "meta" {
			var name, content string

			for i := range child.Attr {
				switch child.Attr[i].Key {
				case "name":
					name = child.Attr[i].Val
				case "content":
					content = child.Attr[i].Val
				}
			}

			// add values to the map
			if name != "" && content != "" {
				if keys == nil {
					keys = make(map[string]string)
				}
				keys[name] = content
			}
		}
	}

	return
}

// finds HTML node
func extractHTMLNode(doc *html.Node, tagName string) (targetNode *html.Node) {
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tagName {
			targetNode = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(doc)
	return
}

// out HTML string with the node and all its content
func nodeToString(node *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, node)
	return buf.String()
}
