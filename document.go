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
	var result []string
	ret := this.walkNode(func(node *html.Node) bool {
		ret := r.FindStringSubmatch(node.Data)
		if len(ret) > 0 && ret[0] == node.Data {
			result = ret[1:]
			return true
		}
		return false
	})
	return result, ret
}

func (this *Document) walkNode(walkFunc func(node *html.Node) bool) bool {
	if walkFunc == nil {
		return false
	}
	stack := []*html.Node{this.root} // init stack
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1] // pop
		if walkFunc(node) {
			return true
		}
		if child := node.FirstChild; child != nil {
			stack = append(stack, child) // push
		}
		if sibling := node.NextSibling; sibling != nil {
			stack = append(stack, sibling) // push
		}
	}
	return false
}
