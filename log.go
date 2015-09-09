package hikaru

import (
	"bytes"
	"github.com/najeira/goutils/nlog"
	"runtime/debug"
)

func (c *Context) logPanic(err interface{}) {
	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()
	c.errorf("%v\n%s", err, stack)
}

func (c *Context) verbosef(format string, args ...interface{}) {
	c.genLogf(nlog.Verbose, format, args...)
}

func (c *Context) debugf(format string, args ...interface{}) {
	c.genLogf(nlog.Debug, format, args...)
}

func (c *Context) infof(format string, args ...interface{}) {
	c.genLogf(nlog.Info, format, args...)
}

func (c *Context) noticef(format string, args ...interface{}) {
	c.genLogf(nlog.Notice, format, args...)
}

func (c *Context) warningf(format string, args ...interface{}) {
	c.genLogf(nlog.Warn, format, args...)
}

func (c *Context) errorf(format string, args ...interface{}) {
	c.genLogf(nlog.Error, format, args...)
}

func (c *Context) criticalf(format string, args ...interface{}) {
	c.genLogf(nlog.Critical, format, args...)
}

func (c *Context) Verbosef(format string, args ...interface{}) {
	c.appLogf(nlog.Verbose, format, args...)
}

func (c *Context) Debugf(format string, args ...interface{}) {
	c.appLogf(nlog.Debug, format, args...)
}

func (c *Context) Infof(format string, args ...interface{}) {
	c.appLogf(nlog.Info, format, args...)
}

func (c *Context) Noticef(format string, args ...interface{}) {
	c.appLogf(nlog.Notice, format, args...)
}

func (c *Context) Warningf(format string, args ...interface{}) {
	c.appLogf(nlog.Warn, format, args...)
}

func (c *Context) Errorf(format string, args ...interface{}) {
	c.appLogf(nlog.Error, format, args...)
}

func (c *Context) Criticalf(format string, args ...interface{}) {
	c.appLogf(nlog.Critical, format, args...)
}
