package hikaru

import (
	"runtime/debug"
)

const (
	LogNo = iota
	LogCritical
	LogError
	LogWarn
	LogInfo
	LogDebug
)

var (
	appLogger Logger
	genLogger Logger
)

type Logger interface {
	V(int) bool
	SetLevel(int)
	Printf(c *Context, level int, format string, args ...interface{})
}

func SetAppLogger(logger Logger) {
	appLogger = logger
}

func SetGenLogger(logger Logger) {
	genLogger = logger
}

func (c *Context) debugf(format string, args ...interface{}) {
	if genLogger != nil {
		genLogger.Printf(c, LogDebug, format, args...)
	}
}

func (c *Context) infof(format string, args ...interface{}) {
	if genLogger != nil {
		genLogger.Printf(c, LogInfo, format, args...)
	}
}

func (c *Context) warningf(format string, args ...interface{}) {
	if genLogger != nil {
		genLogger.Printf(c, LogWarn, format, args...)
	}
}

func (c *Context) errorf(format string, args ...interface{}) {
	if genLogger != nil {
		genLogger.Printf(c, LogError, format, args...)
	}
}

func (c *Context) criticalf(format string, args ...interface{}) {
	if genLogger != nil {
		genLogger.Printf(c, LogCritical, format, args...)
	}
}

func (c *Context) logPanic(err interface{}) {
	c.errorf("%v\n%s", err, string(debug.Stack()))
}

func (c *Context) Debugf(format string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Printf(c, LogDebug, format, args...)
	}
}

func (c *Context) Infof(format string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Printf(c, LogInfo, format, args...)
	}
}

func (c *Context) Warningf(format string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Printf(c, LogWarn, format, args...)
	}
}

func (c *Context) Errorf(format string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Printf(c, LogError, format, args...)
	}
}

func (c *Context) Criticalf(format string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Printf(c, LogCritical, format, args...)
	}
}
