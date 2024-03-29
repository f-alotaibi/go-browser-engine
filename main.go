package main

import (
	"fmt"
	"go-browser-engine/html"
	"go-browser-engine/node"
	"os"
)

func main() {
	buf, err := os.ReadFile("test.html")
	if err != nil {
		panic(err)
	}
	parsedNode := html.Parse(string(buf))
	loopThroughNode(parsedNode)
}

func loopThroughNode(parent node.Node) {
	for i, child := range parent.Children {
		nodeType := child.NodeType
		print(i, ": ")
		if nodeType == node.Text {
			val, ok := child.Value.(string)
			if !ok {
				panic("error in text")
			}
			fmt.Printf("text: %s", val)
		} else if nodeType == node.Element {
			val, ok := child.Value.(node.ElementData)
			if !ok {
				panic("error in text")
			}
			fmt.Printf("Tagname: %s, attrs: %v", val.TagName, val.Attributes)
		}
		print("\n")
		loopThroughNode(child)
	}
}
