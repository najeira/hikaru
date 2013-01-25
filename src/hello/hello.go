package hello

import (
	"fmt"
	"fimika"
)

func init() {
	fimika.AddRoute("/", handleIndex)
	fimika.AddRoute("/<foo>/<bar>", handleIndex)
}

func handleIndex(c *fimika.Context) error {
	fmt.Fprintf(c.Response, "あああ")
	fmt.Fprintf(c.Response, c.Request.Params["foo"][0])
	return nil
}
