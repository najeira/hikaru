package recovery

import (
	"runtime"

	"github.com/najeira/hikaru"
)

func HandlerFunc(h hikaru.HandlerFunc) hikaru.HandlerFunc {
	return func(c *hikaru.Context) {
		defer func() {
			if err := recover(); err != nil {
				handlePanic(c, err)
			}
		}()

		// call handler
		h(c)
	}
}

func handlePanic(c *hikaru.Context, err interface{}) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	c.Errorf("%s\n%s", err, string(buf[:n]))
}
