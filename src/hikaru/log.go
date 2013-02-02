package hikaru

import (
	"appengine"
	"fmt"
)

const (
	LogLevelNo = iota
	LogLevelCritical
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

type logPrinter func(appengine.Context, string, ...interface{})

var levelFuncMap = map[int]logPrinter{
	LogLevelCritical: appengine.Context.Criticalf,
	LogLevelError:    appengine.Context.Errorf,
	LogLevelWarn:     appengine.Context.Warningf,
	LogLevelInfo:     appengine.Context.Infof,
	LogLevelDebug:    appengine.Context.Debugf,
}

func (c *HikaruContext) logPrintf(level int, format string, args ...interface{}) {
	if c.Application.LogLevel >= level {
		f, ok := levelFuncMap[level]
		if ok && f != nil {
			f(c.AppEngineContext(), format, args...)
		}
	}
}

func (c *HikaruContext) Debug(args ...interface{}) {
	c.logPrintf(LogLevelDebug, fmt.Sprint(args...))
}

func (c *HikaruContext) Debugln(args ...interface{}) {
	c.logPrintf(LogLevelDebug, fmt.Sprintln(args...))
}

func (c *HikaruContext) Debugf(format string, args ...interface{}) {
	c.logPrintf(LogLevelDebug, format, args...)
}

func (c *HikaruContext) Info(args ...interface{}) {
	c.logPrintf(LogLevelInfo, fmt.Sprint(args...))
}

func (c *HikaruContext) Infoln(args ...interface{}) {
	c.logPrintf(LogLevelInfo, fmt.Sprintln(args...))
}

func (c *HikaruContext) Infof(format string, args ...interface{}) {
	c.logPrintf(LogLevelInfo, format, args...)
}

func (c *HikaruContext) Warning(args ...interface{}) {
	c.logPrintf(LogLevelWarn, fmt.Sprint(args...))
}

func (c *HikaruContext) Warningln(args ...interface{}) {
	c.logPrintf(LogLevelWarn, fmt.Sprintln(args...))
}

func (c *HikaruContext) Warningf(format string, args ...interface{}) {
	c.logPrintf(LogLevelWarn, format, args...)
}

func (c *HikaruContext) Error(args ...interface{}) {
	c.logPrintf(LogLevelError, fmt.Sprint(args...))
}

func (c *HikaruContext) Errorln(args ...interface{}) {
	c.logPrintf(LogLevelError, fmt.Sprintln(args...))
}

func (c *HikaruContext) Errorf(format string, args ...interface{}) {
	c.logPrintf(LogLevelError, format, args...)
}

func (c *HikaruContext) Critical(args ...interface{}) {
	c.logPrintf(LogLevelCritical, fmt.Sprint(args...))
}

func (c *HikaruContext) Criticalln(args ...interface{}) {
	c.logPrintf(LogLevelCritical, fmt.Sprintln(args...))
}

func (c *HikaruContext) Criticalf(format string, args ...interface{}) {
	c.logPrintf(LogLevelCritical, format, args...)
}
