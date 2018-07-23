package crawler

import (
	"regexp"

	"golang.org/x/net/html"
)

type Document struct {
	root *html.Node
}

func NewDocument(root *html.Node) *Document {
	document := &Document{
		root: root,
	}
	return document
}

func (this *Document) Find(pattern string) ([]string, bool) {
	r := regexp.MustCompile(pattern)

	var stack []*html.Node
	stack = append(stack, this.root) // push
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1] // pop
		if ret := r.FindStringSubmatch(node.Data); ret != nil {
			return ret[1:], true
		}
		if child := node.FirstChild; child != nil {
			stack = append(stack, child) // push
		}
		if sibling := node.NextSibling; sibling != nil {
			stack = append(stack, sibling) // push
		}
	}
	return nil, false
}
