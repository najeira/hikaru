package hello

import (
	"fimika"
)

var (
	app = fimika.NewApplication()
)

func init() {
	app.Debug = true
	app.Route("/", handleIndex)
	app.Route("/<foo>/<bar>", handleIndex)
	app.Start()
}

func handleIndex(c *fimika.Context) *fimika.Result {
	c.LogInfof("こんにちは")
	return c.Text("あああ")
}
