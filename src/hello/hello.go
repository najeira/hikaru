package hello

import (
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

func handleIndex(c *fimika.Context) *fimika.Result {
	return c.Raw("あああ", "text/plain; charset=utf-8")
}
