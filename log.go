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

var (
	applicationLogger Logger
	generalLogger     Logger
)

type Logger interface {
	V(int) bool
	Printf(c *Context, level int, format string, args ...interface{})
}

func init() {
	SetLogger(NewLogger(LogInfo))
	SetGeneralLogger(NewLogger(LogWarn))
}

func SetLogger(logger Logger) {
	applicationLogger = logger
}

func SetGeneralLogger(logger Logger) {
	generalLogger = logger
}

func logGenf(c *Context, level int, format string, args ...interface{}) {
	if generalLogger != nil && generalLogger.V(level) {
		generalLogger.Printf(c, level, format, args...)
	}
}

func logAppf(c *Context, level int, format string, args ...interface{}) {
	if applicationLogger != nil && applicationLogger.V(level) {
		applicationLogger.Printf(c, level, format, args...)
	}
}

func (c *Context) tracef(format string, args ...interface{}) {
	logGenf(c, LogTrace, format, args...)
}

func (c *Context) debugf(format string, args ...interface{}) {
	logGenf(c, LogDebug, format, args...)
}

func (c *Context) infof(format string, args ...interface{}) {
	logGenf(c, LogInfo, format, args...)
}

func (c *Context) warningf(format string, args ...interface{}) {
	logGenf(c, LogWarn, format, args...)
}

func (c *Context) errorf(format string, args ...interface{}) {
	logGenf(c, LogError, format, args...)
}

func (c *Context) criticalf(format string, args ...interface{}) {
	logGenf(c, LogCritical, format, args...)
}

func (c *Context) Tracef(format string, args ...interface{}) {
	logAppf(c, LogTrace, format, args...)
}

func (c *Context) Debugf(format string, args ...interface{}) {
	logAppf(c, LogDebug, format, args...)
}

func (c *Context) Infof(format string, args ...interface{}) {
	logAppf(c, LogInfo, format, args...)
}

func (c *Context) Warningf(format string, args ...interface{}) {
	logAppf(c, LogWarn, format, args...)
}

func (c *Context) Errorf(format string, args ...interface{}) {
	logAppf(c, LogError, format, args...)
}

func (c *Context) Criticalf(format string, args ...interface{}) {
	logAppf(c, LogCritical, format, args...)
}
