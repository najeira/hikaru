package hikaru

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	request *http.Request
	status  int
	size    int
}

var (
	_ http.ResponseWriter = (*Response)(nil)
	_ http.Hijacker       = (*Response)(nil)
	_ http.Flusher        = (*Response)(nil)
	_ http.CloseNotifier  = (*Response)(nil)
)

func (r *Response) init(w http.ResponseWriter, req *http.Request) {
	r.ResponseWriter = w
	r.request = req
	r.status = http.StatusOK
	r.size = -1
}

// Returns response headers.
func (r *Response) Header() http.Header {
	return r.ResponseWriter.Header()
}

// GetHeader gets a response header.
func (r *Response) GetHeader(key string) string {
	return r.Header().Get(key)
}

// SetHeader sets a response header.
func (r *Response) SetHeader(key, value string) {
	r.Header().Set(key, value)
}

// Adds a response header.
func (r *Response) AddHeader(key, value string) {
	r.Header().Add(key, value)
}

// Adds a cookie header.
func (r *Response) SetCookie(cookie *http.Cookie) {
	r.Header().Set("Set-Cookie", cookie.String())
}

func (r *Response) Status() int {
	return r.status
}

func (r *Response) Size() int {
	return r.size
}

func (r *Response) Written() bool {
	return r.size >= 0
}

func (r *Response) WriteHeader(code int) {
	if code > 0 && r.status != code {
		r.status = code
	}
}

func (r *Response) WriteHeaderAndSend(code int) {
	r.WriteHeader(code)
	r.writeHeaderIfNotSent()
}

func (r *Response) writeHeaderIfNotSent() {
	if !r.Written() {
		r.size = 0
		r.ResponseWriter.WriteHeader(r.status)
	}
}

func (r *Response) Write(msg []byte) (int, error) {
	r.writeHeaderIfNotSent()
	n, err := r.ResponseWriter.Write(msg)
	r.size += n
	return n, err
}

// Writes raw bytes and content type.
func (r *Response) Raw(body []byte, contentType string) (int, error) {
	if contentType != "" {
		r.SetHeader("Content-Type", contentType)
	}
	return r.Write(body)
}

// Writes a text string.
// The content type should be "text/plain; charset=utf-8".
func (r *Response) Text(body string) (int, error) {
	r.SetHeader("Content-Type", "text/plain; charset=utf-8")
	return io.WriteString(r, body)
}

func (r *Response) Json(value interface{}) error {
	r.SetHeader("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(r)
	if err := e.Encode(value); err != nil {
		return err
	}
	return nil
}

// Sets response to HTTP 3xx.
func (r *Response) Redirect(path string, code int) {
	r.SetHeader("Location", path)
	http.Redirect(r, r.request, path, code)
	r.writeHeaderIfNotSent()
}

// Sets response to HTTP 304 Not Modified.
func (r *Response) NotModified() {
	r.WriteHeaderAndSend(http.StatusNotModified)
}

// Sets response to HTTP 401 Unauthorized.
func (r *Response) Unauthorized() {
	r.WriteHeaderAndSend(http.StatusUnauthorized)
}

// Sets response to HTTP 403 Forbidden.
func (r *Response) Forbidden() {
	r.WriteHeaderAndSend(http.StatusForbidden)
}

// Sets response to HTTP 404 Not Found.
func (r *Response) NotFound() {
	r.WriteHeaderAndSend(http.StatusNotFound)
}

// Sets response to HTTP 500 Internal Server Error.
func (r *Response) Fail(err interface{}) {
	r.WriteHeaderAndSend(http.StatusInternalServerError)
}

// Hijack lets the caller take over the connection.
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if r.size < 0 {
		r.size = 0
	}
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify returns a channel that receives a single value
// when the client connection has gone away.
func (r *Response) CloseNotify() <-chan bool {
	return r.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush sends any buffered data to the client.
func (r *Response) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}
