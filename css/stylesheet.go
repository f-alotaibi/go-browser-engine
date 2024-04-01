package css

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type Rule struct {
	selectors    []Selector
	declarations []Declaration
}

type Stylesheet struct {
	rules []Rule
}

type Selector struct {
	selectorType SelectorType
	value        interface{}
}
type SelectorType uint16

const (
	Simple SelectorType = iota
)

type SimpleSelector struct {
	tagName string // Optional
	id      string // Optional
	class   []string
}

type Declaration struct {
	name  string
	value Value
}

type ValueType uint16

const (
	Keyword ValueType = iota
	Length
	ColorValue
)

type Value struct {
	valueType ValueType
	value     interface{}
}

type Unit uint16

const (
	Px Unit = iota
)

type Color struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type parser struct {
	position uint16
	input    string
}

// returns the next char? it should be current
func (parser *parser) NextChar() string {
	return strings.Split(parser.input, "")[parser.position]
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

func (parser *parser) ConsumeComment() {
	if !parser.StartsWith("<!--") {
		return
	}
	text := ""
	parser.ConsumeWhile(func(s string) bool {
		text += s
		return !strings.HasSuffix(text, "-->")
	})
	if !parser.EOF() {
		parser.ConsumeChar()
	}
}

func (parser *parser) ParseIdentifier() string {
	return parser.ConsumeWhile(validIdentifierChar)
}

func (parser *parser) ParseSimpleSelector() SimpleSelector {
	selector := SimpleSelector{
		class: make([]string, 0),
	}
	for !parser.EOF() {
		nextChar := parser.NextChar()
		if nextChar == "#" {
			parser.ConsumeChar()
			selector.id = parser.ParseIdentifier()
		} else if nextChar == "." {
			parser.ConsumeChar()
			selector.class = append(selector.class, parser.ParseIdentifier())
		} else if nextChar == "*" {
			// universal
			parser.ConsumeChar()
		} else if validIdentifierChar(nextChar) {
			selector.tagName = parser.ParseIdentifier()
		}
	}
	return selector
}

func (parser *parser) ParseRule() Rule {
	return Rule{
		selectors:    parser.ParseSelectors(),
		declarations: parser.ParseDeclarations(),
	}
}

func (parser *parser) ParseDeclarations() []Declaration {
	if parser.ConsumeChar() != "{" {
		panic("Cannot declare because char is not {")
	}
	declarations := make([]Declaration, 0)
	for {
		parser.ConsumeWhitespace()
		if parser.NextChar() == "}" {
			parser.ConsumeChar()
			break
		}
		declarations = append(declarations, parser.ParseDeclaration())
	}
	return declarations
}

func (parser *parser) ParseDeclaration() Declaration {
	propertyName := parser.ParseIdentifier()
	parser.ConsumeWhitespace()
	if parser.ConsumeChar() != ":" {
		panic(": error")
	}
	parser.ConsumeWhitespace()
	value := parser.ParseValue()
	parser.ConsumeWhitespace()
	if parser.ConsumeChar() != ";" {
		panic("; error")
	}

	return Declaration{
		name:  propertyName,
		value: value,
	}
}

func (parser *parser) ParseValue() Value {
	switch parser.NextChar() {
	case "0":
	case "1":
	case "2":
	case "3":
	case "4":
	case "5":
	case "6":
	case "7":
	case "8":
	case "9":
		return parser.ParseLength()
	case "#":
		return parser.ParseColor()
	}
	return Value{
		valueType: Keyword,
		value:     parser.ParseIdentifier(),
	}
}

type LengthValue struct {
	length float32
	unit   Unit
}

func (parser *parser) ParseLength() Value {
	return Value{
		valueType: Length,
		value: LengthValue{
			length: parser.ParseFloat(),
			unit:   parser.ParseUnit(),
		},
	}
}

func (parser *parser) ParseFloat() float32 {
	s := parser.ConsumeWhile(func(s string) bool {
		matched, _ := regexp.Match("[0-9.]", []byte(s))
		return matched
	})
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

func (parser *parser) ParseUnit() Unit {
	identifier := strings.ToLower(parser.ParseIdentifier())
	if identifier == "px" {
		return Px
	} else {
		panic(fmt.Sprintf("unrecgnized unit %s", identifier))
	}
}

func (parser *parser) ParseColor() Value {
	if parser.ConsumeChar() != "#" {
		panic("Invalid charcater color")
	}
	return Value{
		valueType: ColorValue,
		value: Color{
			r: parser.ParseHexPair(),
			g: parser.ParseHexPair(),
			b: parser.ParseHexPair(),
			a: 255,
		},
	}
}

func (parser *parser) ParseHexPair() uint8 {
	i := 0
	s := parser.ConsumeWhile(func(s string) bool {
		i++
		return i <= 2
	})
	value, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

func (parser *parser) ParseSelectors() []Selector {
	selectors := make([]Selector, 0)
	for {
		selectorValue := parser.ParseSimpleSelector()
		selectors = append(selectors, Selector{
			selectorType: Simple,
			value:        selectorValue,
		})
		parser.ConsumeWhitespace()
		nextChar := parser.NextChar()
		if nextChar == "," {
			parser.ConsumeChar()
			parser.ConsumeWhitespace()
		} else if nextChar == "{" {
			break
		} else {
			panic(fmt.Sprintf("Unexpected character %s", nextChar))
		}
	}

	slices.SortFunc(selectors, func(a, b Selector) int {
		return cmp.Compare((a.Specificity().A + a.Specificity().B + a.Specificity().C), (b.Specificity().A + b.Specificity().B + b.Specificity().C))
	})
	return selectors
}

func validIdentifierChar(s string) bool {
	matched, _ := regexp.Match("[A-Za-z0-9-_]", []byte(s))
	return matched
}

type Specificity struct {
	A, B, C int
}

func (selector *Selector) Specificity() Specificity {
	simple := *selector
	value, ok := simple.value.(SimpleSelector)
	if !ok {
		return Specificity{}
	}
	a := len(value.id)
	b := len(value.class)
	c := len(value.tagName)
	return Specificity{
		A: a,
		B: b,
		C: c,
	}
}
