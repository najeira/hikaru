package hikaru

import (
	"net/http"
	"sync"
	"time"
)

type timeoutHandler struct {
	timeout time.Duration
}

func TimeoutHandler(t time.Duration) HandlerFunc {
	h := &timeoutHandler{t}
	return h.handle
}

func (h *timeoutHandler) handle(c *Context) {
	errCh := make(chan interface{}, 1)

	// Hijack ResponseWriter
	tr := &timeoutResponse{ResponseWriter: c.ResponseWriter}
	c.ResponseWriter = tr

	// Start handlers on the new goroutine.
	go func() {
		// Fetch the handler's panic if exists.
		defer func() {
			if err := recover(); err != nil {
				// Send error.
				errCh <- err
			} else {
				// No error.
				close(errCh)
			}
		}()
		c.Next()
	}()

	// Wait handlers or timed out.
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-time.After(h.timeout):
		tr.WriteHeader(http.StatusServiceUnavailable)
	}
}

type timeoutResponse struct {
	http.ResponseWriter

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
}

func (r *timeoutResponse) Write(p []byte) (int, error) {
	r.mu.Lock()
	timedOut := r.timedOut
	r.mu.Unlock()
	if timedOut {
		return 0, http.ErrHandlerTimeout
	}
	return r.ResponseWriter.Write(p)
}

func (r *timeoutResponse) WriteHeader(code int) {
	r.mu.Lock()
	done := r.timedOut || r.wroteHeader
	r.wroteHeader = true
	r.mu.Unlock()
	if !done {
		r.ResponseWriter.WriteHeader(code)
	}
}
