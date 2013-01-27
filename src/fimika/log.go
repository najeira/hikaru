package fimika

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

func (c *Context) logPrintf(level int, format string, args ...interface{}) {
	if c.Application.LogLevel >= level {
		f, ok := levelFuncMap[level]
		if ok && f != nil {
			f(c.AppEngineContext, format, args...)
		}
	}
}

func (c *Context) LogDebug(args ...interface{}) {
	c.logPrintf(LogLevelDebug, fmt.Sprint(args...))
}

func (c *Context) LogDebugln(args ...interface{}) {
	c.logPrintf(LogLevelDebug, fmt.Sprintln(args...))
}

func (c *Context) LogDebugf(format string, args ...interface{}) {
	c.logPrintf(LogLevelDebug, format, args...)
}

func (c *Context) LogInfo(args ...interface{}) {
	c.logPrintf(LogLevelInfo, fmt.Sprint(args...))
}

func (c *Context) LogInfoln(args ...interface{}) {
	c.logPrintf(LogLevelInfo, fmt.Sprintln(args...))
}

func (c *Context) LogInfof(format string, args ...interface{}) {
	c.logPrintf(LogLevelInfo, format, args...)
}

func (c *Context) LogWarn(args ...interface{}) {
	c.logPrintf(LogLevelWarn, fmt.Sprint(args...))
}

func (c *Context) LogWarnln(args ...interface{}) {
	c.logPrintf(LogLevelWarn, fmt.Sprintln(args...))
}

func (c *Context) LogWarnf(format string, args ...interface{}) {
	c.logPrintf(LogLevelWarn, format, args...)
}

func (c *Context) LogError(args ...interface{}) {
	c.logPrintf(LogLevelError, fmt.Sprint(args...))
}

func (c *Context) LogErrorln(args ...interface{}) {
	c.logPrintf(LogLevelError, fmt.Sprintln(args...))
}

func (c *Context) LogErrorf(format string, args ...interface{}) {
	c.logPrintf(LogLevelError, format, args...)
}

func (c *Context) LogCritical(args ...interface{}) {
	c.logPrintf(LogLevelCritical, fmt.Sprint(args...))
}

func (c *Context) LogCriticalln(args ...interface{}) {
	c.logPrintf(LogLevelCritical, fmt.Sprintln(args...))
}

func (c *Context) LogCriticalf(format string, args ...interface{}) {
	c.logPrintf(LogLevelCritical, format, args...)
}
