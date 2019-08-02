package main

import (
	"os"

	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	p.Task("server", nil, func(c *do.Context) {
		target := os.Getenv("SERVER")
		switch target {
		case "WebScoket", "websocket":
			c.Start("websocket/main.go", do.M{"$in": "./websocket/"})
		default:
			c.Start("main.go", do.M{"$in": "./"})
		}
	}).Src("*.go", "**/*.go").
		Debounce(3000)
}

func main() {
	do.Godo(tasks)
}
