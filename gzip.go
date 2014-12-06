package hikaru

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponse struct {
	response   http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// gzipResponse should be http.ResponseWriter
var _ http.ResponseWriter = (*gzipResponse)(nil)

var GzipAcceptableContentTypes []string = []string{"text/", "application/json"}

func (r *gzipResponse) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *gzipResponse) WriteHeader(code int) {
	r.statusCode = code
}

func (r *gzipResponse) Header() http.Header {
	return r.response.Header()
}

func (r *gzipResponse) flush() {
	// Get wether the content type is supported
	headers := r.response.Header()
	ct := headers.Get("Content-Type")
	gzipable := false
	for _, act := range GzipAcceptableContentTypes {
		if strings.HasPrefix(ct, act) {
			gzipable = true
			break
		}
	}
	if gzipable {
		// Set the appropriate gzip headers.
		headers.Set("Content-Encoding", "gzip")
		headers.Set("Vary", "Accept-Encoding")

		// Delete the content length after we know we have been written to.
		headers.Del("Content-Length")

		// Compress and Write
		gw := gzip.NewWriter(r.response)
		defer gw.Close()
		gw.Write(r.body.Bytes())
	} else {
		// Write original response directly
		r.response.WriteHeader(r.statusCode)
		r.response.Write(r.body.Bytes())
	}
}

func GzipHandlerFunc(c *Context) {
	// Skip compression if the client doesn't accept gzip encoding.
	if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
		return
	}

	// Set a gzipResponse to the Context to wrap response
	res := &gzipResponse{
		response:   c.ResponseWriter,
		statusCode: http.StatusOK,
	}
	c.ResponseWriter = res

	// Call the next handler.
	c.Next()

	// Flush the response
	res.flush()
}
