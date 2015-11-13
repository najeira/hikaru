// +build appengine

package hikaru

import (
	"appengine"
)

var logLevelAppEngineLoggerMap = map[int](func(appengine.Context, string, ...interface{})){
	LogDebug:    appengine.Context.Debugf,
	LogInfo:     appengine.Context.Infof,
	LogWarn:     appengine.Context.Warningf,
	LogError:    appengine.Context.Errorf,
	LogCritical: appengine.Context.Criticalf,
}

type appengineLogger struct {
	level int
}

func newDefaultLogger(level int) Logger {
	return &appengineLogger{level: level}
}

func (l *appengineLogger) V(level int) bool {
	return l.level <= level && level > LogNo
}

func (l *appengineLogger) SetLevel(level int) {
	l.level = level
}

func (l *appengineLogger) Printf(c *Context, level int, format string, args ...interface{}) {
	if l.V(level) {
		if f, ok := logLevelAppEngineLoggerMap[level]; ok {
			f(c.AC(), format, args...)
		}
	}
}

func (c *Context) initEnv() {
	c.appengineContext = appengine.NewContext(c.Request)
}

func (c *Context) AC() appengine.Context {
	return c.appengineContext.(appengine.Context)
}
