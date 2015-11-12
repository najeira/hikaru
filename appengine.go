// +build appengine

package hikaru

import (
	"appengine"
	"net/http"
)

var logLevelAppEngineLoggerMap = map[int](func(appengine.Context, string, ...interface{})){
	LogDebug:    appengine.Context.Debugf,
	LogInfo:     appengine.Context.Infof,
	LogWarn:     appengine.Context.Warningf,
	LogError:    appengine.Context.Errorf,
	LogCritical: appengine.Context.Criticalf,
}

func init() {
	appLogger = &appengineLogger{level: LogDebug}
	genLogger = &appengineLogger{level: LogDebug}
}

type appengineLogger struct {
	level int
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

type envContext struct {
	appengine.Context
}

func (c *envContext) init(r *http.Request) {
	if r != nil {
		c.Context = appengine.NewContext(r)
	} else {
		c.Context = nil
	}
}

func (c *Context) AC() appengine.Context {
	return c.envContext.Context
}
