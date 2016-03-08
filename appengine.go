// +build appengine

package hikaru

import (
	"appengine"
)

var logLevelAppEngineLoggerMap = map[int](func(appengine.Context, string, ...interface{})){
	LogTrace:    appengine.Context.Debugf,
	LogDebug:    appengine.Context.Debugf,
	LogInfo:     appengine.Context.Infof,
	LogWarn:     appengine.Context.Warningf,
	LogError:    appengine.Context.Errorf,
	LogCritical: appengine.Context.Criticalf,
}

type appengineLogger struct {
	level int
}

func NewLogger(level int) Logger {
	return &appengineLogger{level: level}
}

func (l *appengineLogger) V(level int) bool {
	return l.level >= level && level > LogNo
}

func (l *appengineLogger) Printf(c *Context, level int, format string, args ...interface{}) {
	if l.V(level) {
		if f, ok := logLevelAppEngineLoggerMap[level]; ok {
			f(appengine.NewContext(c.Request), format, args...)
		}
	}
}
