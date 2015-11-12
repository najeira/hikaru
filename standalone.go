// +build !appengine

package hikaru

import (
	"fmt"
	"log"
	"net/http"
)

var logLevelNameMap = map[int]string{
	LogCritical: "CRITICAL",
	LogError:    "ERROR",
	LogWarn:     "WARN",
	LogInfo:     "INFO",
	LogDebug:    "DEBUG",
}

func init() {
	appLogger = &defaultLogger{level: LogInfo}
	genLogger = &defaultLogger{level: LogWarn}
}

type defaultLogger struct {
	level int
}

func (l *defaultLogger) V(level int) bool {
	return l.level <= level && level > LogNo
}

func (l *defaultLogger) SetLevel(level int) {
	l.level = level
}

func (l *defaultLogger) Printf(c *Context, level int, format string, args ...interface{}) {
	if l.V(level) {
		if name, ok := logLevelNameMap[level]; ok {
			format2 := fmt.Sprintf("[%s] %s", name, format)
			log.Printf(format2, args...)
		} else {
			log.Printf(format, args...)
		}
	}
}

type envContext struct {
}

func (c *envContext) init(r *http.Request) {
	// nothing for standalone environment
}
