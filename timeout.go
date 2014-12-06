package hikaru

import (
	"net/http"
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
	errCh := make(chan interface{}, 0)

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
			c.WriteHeader(http.StatusInternalServerError)
		}
	case <-time.After(h.timeout):
		c.WriteHeader(http.StatusServiceUnavailable)
	}
}
