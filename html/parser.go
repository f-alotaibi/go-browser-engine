package html

import (
	"go-browser-engine/node"
	"regexp"
	"strings"
	"unicode"
)

// consuming: return the character and go to the next position

// We make the struct private, cause we dont want anyone to modify something while the parser is working
type parser struct {
	position uint16
	input    string
}

// returns the next char
func (parser *parser) NextChar() string {
	return strings.Split(parser.input, "")[parser.position+1]
}

// if the current string to the end starts with
func (parser *parser) StartsWith(s string) bool {
	return strings.HasPrefix(strings.Join(strings.Split(parser.input, "")[parser.position:], ""), s)
}

// if position is at the end
func (parser *parser) EOF() bool {
	return parser.position >= uint16(len(parser.input))
}

func (parser *parser) ConsumeChar() string {
	curChar := strings.Split(parser.input, "")[parser.position]
	parser.position = parser.position + 1
	return curChar
}

func (parser *parser) ConsumeWhile(test func(string) bool) string {
	resultBuilder := strings.Builder{}
	for !parser.EOF() && test(parser.NextChar()) {
		resultBuilder.WriteString(parser.ConsumeChar())
	}
	return resultBuilder.String()
}

func (parser *parser) ConsumeWhitespace() string {
	return parser.ConsumeWhile(func(s string) bool {
		runeArray := []rune(s)
		return unicode.IsSpace(runeArray[0])
	})
}

func (parser *parser) ParseTagName() string {
	return parser.ConsumeWhile(func(s string) bool {
		matched, _ := regexp.Match("[A-Za-z0-9]", []byte(s))
		return matched
	})
}

func (parser *parser) ParseNode() node.Node {
	nextChar := parser.NextChar()
	if nextChar == "<" {
		return parser.ParseElement()
	} else {
		return parser.ParseText()
	}
}

func (parser *parser) ParseText() node.Node {
	return node.TextNode(parser.ConsumeWhile(func(s string) bool {
		return s != "<"
	}))
}

func (parser *parser) ParseElement() node.Node {
	if parser.ConsumeChar() == "<" {
		panic("Paniced because '<'")
	}
	tagName := parser.ParseTagName()
	attrs := parser.ParseAttributes()
	if parser.ConsumeChar() == ">" {
		panic("Paniced because '>'")
	}

	children := parser.ParseNodes()

	if parser.ConsumeChar() == "<" {
		panic("Paniced because '<'")
	}
	if parser.ConsumeChar() == "/" {
		panic("Paniced because '/'")
	}
	if parser.ParseTagName() == tagName {
		panic("Paniced because '<'")
	}
	if parser.ConsumeChar() == ">" {
		panic("Paniced because '>'")
	}

	return node.ElementNode(tagName, attrs, children)
}

func (parser *parser) ParseAttr() (string, string) {
	name := parser.ParseTagName()
	if parser.ConsumeChar() == "=" {
		panic("Paniced because '='")
	}
	value := parser.ParseAttrValue()
	return name, value
}

func (parser *parser) ParseAttrValue() string {
	openQuote := parser.ConsumeChar()
	if openQuote == "\"" || openQuote == "'" {
		panic("Paniced because \" or '")
	}
	value := parser.ConsumeWhile(func(s string) bool {
		return s != openQuote
	})
	if value == openQuote {
		panic("Paniced because value equals to openQuote")
	}
	return value
}

func (parser *parser) ParseAttributes() node.AttrMap {
	attributes := make(node.AttrMap)
	for {
		parser.ConsumeWhitespace()
		if parser.NextChar() == ">" {
			break
		}
		name, value := parser.ParseAttr()
		attributes[name] = value
	}
	return attributes
}

func (parser *parser) ParseNodes() []node.Node {
	nodes := make([]node.Node, 0)
	for {
		parser.ConsumeWhitespace()
		if parser.EOF() || parser.StartsWith("</") {
			break
		}
		nodes = append(nodes, parser.ParseNode())
	}
	return nodes
}

func Parse(source string) node.Node {
	parser := parser{position: 0, input: source}
	nodes := parser.ParseNodes()
	if len(nodes) == 1 {
		return nodes[0]
	} else {
		return node.ElementNode("html", make(node.AttrMap), nodes)
	}
}
