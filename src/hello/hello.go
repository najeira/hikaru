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
	return nil
}
