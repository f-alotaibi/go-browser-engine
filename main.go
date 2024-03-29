package main

import "go-browser-engine/html"

func main() {
	node := html.Parse("<body>Hello, world!</body>")
	println(node.Value)
}
