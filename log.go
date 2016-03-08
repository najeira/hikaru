package hikaru

const (
	LogNo = iota
	LogCritical
	LogError
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

type Logger interface {
	V(int) bool
	Printf(c *Context, level int, format string, args ...interface{})
}

type logger struct {
	app Logger
	gen Logger
}

func (l *logger) genf(c *Context, level int, format string, args ...interface{}) {
	if l.gen != nil && l.gen.V(level) {
		l.gen.Printf(c, level, format, args...)
	}
}

func (l *logger) genv(c *Context, level int, value interface{}) {
	if l.gen != nil && l.gen.V(level) {
		l.gen.Printf(c, level, "%v", value)
	}
}

func (l *logger) appf(c *Context, level int, format string, args ...interface{}) {
	if l.app != nil && l.app.V(level) {
		l.app.Printf(c, level, format, args...)
	}
}

func (l *logger) appv(c *Context, level int, value interface{}) {
	if l.app != nil && l.app.V(level) {
		l.app.Printf(c, level, "%v", value)
	}
}

func (c *Context) tracef(format string, args ...interface{}) {
	c.logger.genf(c, LogTrace, format, args...)
}

func (c *Context) debugf(format string, args ...interface{}) {
	c.logger.genf(c, LogDebug, format, args...)
}

func (c *Context) infof(format string, args ...interface{}) {
	c.logger.genf(c, LogInfo, format, args...)
}

func (c *Context) warningf(format string, args ...interface{}) {
	c.logger.genf(c, LogWarn, format, args...)
}

func (c *Context) errorf(format string, args ...interface{}) {
	c.logger.genf(c, LogError, format, args...)
}

func (c *Context) criticalf(format string, args ...interface{}) {
	c.logger.genf(c, LogCritical, format, args...)
}

func (c *Context) Tracef(format string, args ...interface{}) {
	c.logger.appf(c, LogTrace, format, args...)
}

func (c *Context) Debugf(format string, args ...interface{}) {
	c.logger.appf(c, LogDebug, format, args...)
}

func (c *Context) Infof(format string, args ...interface{}) {
	c.logger.appf(c, LogInfo, format, args...)
}

func (c *Context) Warningf(format string, args ...interface{}) {
	c.logger.appf(c, LogWarn, format, args...)
}

func (c *Context) Errorf(format string, args ...interface{}) {
	c.logger.appf(c, LogError, format, args...)
}

func (c *Context) Criticalf(format string, args ...interface{}) {
	c.logger.appf(c, LogCritical, format, args...)
}
