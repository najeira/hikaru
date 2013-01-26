package hello

import (
	"fmt"
	"fimika"
)

var (
	app = fimika.NewApplication()
)

func init() {
	app.AddRoute("/", handleIndex)
	app.AddRoute("/<foo>/<bar>", handleIndex)
	app.Start()
}

func handleIndex(c *fimika.Context) {
	fmt.Fprintf(c.ResponseWriter, "あああ")
	fmt.Fprintf(c.ResponseWriter, c.Get("foo"))
}
