package hikaru

import (
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

type Logger interface {
	SetLevel(level int)
	Write(level int, message string)
	Flush()
}

type LogLogger struct {
	level int
}

func (c *Context) logPrint(level int, message string) {
}

func (c *Context) Debug(args ...interface{}) {
	c.logPrint(LogLevelDebug, fmt.Sprint(args...))
}

func (c *Context) Debugln(args ...interface{}) {
	c.logPrint(LogLevelDebug, fmt.Sprintln(args...))
}

func (c *Context) Debugf(format string, args ...interface{}) {
	c.logPrint(LogLevelDebug, fmt.Sprintf(format, args...))
}

func (c *Context) Info(args ...interface{}) {
	c.logPrint(LogLevelInfo, fmt.Sprint(args...))
}

func (c *Context) Infoln(args ...interface{}) {
	c.logPrint(LogLevelInfo, fmt.Sprintln(args...))
}

func (c *Context) Infof(format string, args ...interface{}) {
	c.logPrint(LogLevelInfo, fmt.Sprintf(format, args...))
}

func (c *Context) Warning(args ...interface{}) {
	c.logPrint(LogLevelWarn, fmt.Sprint(args...))
}

func (c *Context) Warningln(args ...interface{}) {
	c.logPrint(LogLevelWarn, fmt.Sprintln(args...))
}

func (c *Context) Warningf(format string, args ...interface{}) {
	c.logPrint(LogLevelWarn, fmt.Sprintf(format, args...))
}

func (c *Context) Error(args ...interface{}) {
	c.logPrint(LogLevelError, fmt.Sprint(args...))
}

func (c *Context) Errorln(args ...interface{}) {
	c.logPrint(LogLevelError, fmt.Sprintln(args...))
}

func (c *Context) Errorf(format string, args ...interface{}) {
	c.logPrint(LogLevelError, fmt.Sprintf(format, args...))
}

func (c *Context) Critical(args ...interface{}) {
	c.logPrint(LogLevelCritical, fmt.Sprint(args...))
}

func (c *Context) Criticalln(args ...interface{}) {
	c.logPrint(LogLevelCritical, fmt.Sprintln(args...))
}

func (c *Context) Criticalf(format string, args ...interface{}) {
	c.logPrint(LogLevelCritical, fmt.Sprintf(format, args...))
}
