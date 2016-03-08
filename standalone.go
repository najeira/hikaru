// +build !appengine

package hikaru

import (
	"log"
)

var logLevelNameMap = map[int]string{
	LogTrace:    "[TRACE]",
	LogDebug:    "[DEBUG]",
	LogInfo:     "[INFO] ",
	LogWarn:     "[WARN] ",
	LogError:    "[ERROR]",
	LogCritical: "[CRIT] ",
}

type defaultLogger struct {
	level int
}

func NewLogger(level int) Logger {
	return &defaultLogger{level: level}
}

func (l *defaultLogger) V(level int) bool {
	return l.level >= level && level > LogNo
}

func (l *defaultLogger) Printf(c *Context, level int, format string, args ...interface{}) {
	if l.V(level) {
		if name, ok := logLevelNameMap[level]; ok {
			log.Printf(name+format, args...)
		} else {
			log.Printf(format, args...)
		}
	}
}
