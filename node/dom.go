package node

type Node struct {
	Children []Node
	NodeType NodeType
	Value    interface{}
}

type NodeType int

const (
	Text NodeType = iota
	Element
)

type AttrMap map[string]string

type ElementData struct {
	TagName    string
	Attributes AttrMap
}

func TextNode(data string) Node {
	return Node{
		Children: make([]Node, 0),
		NodeType: Text,
		Value:    data,
	}
}

func ElementNode(name string, attrs AttrMap, children []Node) Node {
	return Node{
		Children: children,
		NodeType: Element,
		Value: ElementData{
			TagName:    name,
			Attributes: attrs,
		},
	}
}
