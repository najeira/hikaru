// from https://raw.githubusercontent.com/phyber/negroni-gzip

package hikaru

import (
	"compress/gzip"
	"strings"
)

type gzipHandler struct {
}

func Gzip() HandlerFunc {
	h := &gzipHandler{}
	return h.handle
}

var gzipAcceptableContentTypes []string = []string{"text/", "application/json"}

func (h *gzipHandler) handle(c *Context) {
	// Skip compression if the client doesn't accept gzip encoding.
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		return
	}

	// Call the next handler.
	c.Next()

	if c.body == nil || c.body.Len() <= 0 {
		return
	}

	// Skip compression if the content type is not supported
	headers := c.Header()
	ct := headers.Get("Content-Type")
	gzipable := false
	for _, act := range gzipAcceptableContentTypes {
		if strings.HasPrefix(ct, act) {
			gzipable = true
			break
		}
	}
	if !gzipable {
		return
	}

	// Set the appropriate gzip headers.
	headers.Set("Content-Encoding", "gzip")
	headers.Set("Vary", "Accept-Encoding")

	// Delete the content length after we know we have been written to.
	headers.Del("Content-Length")

	gw := gzip.NewWriter(c.res)
	defer gw.Close()
	c.writeToWriter(gw)
}
