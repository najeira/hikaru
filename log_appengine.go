// +build appengine

package hikaru

import (
	"appengine"
)

type logPrinter func(appengine.Context, string, ...interface{})

var levelFuncMap = map[int]logPrinter{
	LogLevelCritical: appengine.Context.Criticalf,
	LogLevelError:    appengine.Context.Errorf,
	LogLevelWarn:     appengine.Context.Warningf,
	LogLevelInfo:     appengine.Context.Infof,
	LogLevelDebug:    appengine.Context.Debugf,
}

type AppEngineLogger struct {
	level   int
	context appengine.Context
}

func (l *AppEngineLogger) SetLevel(level int) {
	l.level = level
}

func (l *AppEngineLogger) Write(level int, message string) {
	if level > l.level {
		return
	}
	f, ok := levelFuncMap[l.level]
	if ok {
		f(l.context, "%s", message)
	}
}

func (l *AppEngineLogger) Flush() {

}
