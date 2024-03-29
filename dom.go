package main

type Node struct {
	children []Node
	nodeType NodeType
	value    interface{}
}

type NodeType int

const (
	Text NodeType = iota
	Element
)

type AttrMap map[string]string

type ElementData struct {
	tagName    string
	attributes AttrMap
}

func text(data string) Node {
	return Node{
		children: make([]Node, 0),
		nodeType: Text,
		value:    data,
	}
}

func elem(name string, attrs AttrMap, children []Node) Node {
	return Node{
		children: children,
		nodeType: Element,
		value: ElementData{
			tagName:    name,
			attributes: attrs,
		},
	}
}
