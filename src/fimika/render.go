package fimika

import (
	"io"
)

type Renderer interface {
	Render(file string, viewData map[string]interface{}, w io.Writer)
	Ext() string
}


