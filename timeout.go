package hikaru

import (
	"net/http"
	"time"
)

type timeoutHandler struct {
	timeout time.Duration
	msg     []byte
}

func Timeout(t time.Duration, msg string) HandlerFunc {
	h := &timeoutHandler{t, []byte(msg)}
	return h.handle
}

func (h *timeoutHandler) handle(c *Context) {
	done := make(chan bool, 1)
	go func() {
		h.nextWithRecover(c)
		done <- true
	}()
	select {
	case <-done:
		return
	case <-time.After(h.timeout):
		c.write(http.StatusServiceUnavailable, h.msg)
	}
}

func (h *timeoutHandler) nextWithRecover(c *Context) {
	defer c.recover()
	c.Next()
}
