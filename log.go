package hikaru

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

const (
	LogLevelNo = iota
	LogLevelCritical
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

var logLevelName = map[int]string{
	LogLevelCritical: "CRITICAL",
	LogLevelError:    "ERROR",
	LogLevelWarn:     "WARN",
	LogLevelInfo:     "INFO",
	LogLevelDebug:    "DEBUG",
}

type Logger interface {
	SetLevel(level int)
	Write(level int, message []byte)
	Flush()
}

type BufioLogger struct {
	level int
	out   *bufio.Writer
	mu    sync.Mutex
}

func NewBufioLogger(out *bufio.Writer) *BufioLogger {
	return &BufioLogger{
		level: LogLevelDebug,
		out:   out,
	}
}

func NewStdoutLogger() Logger {
	return NewBufioLogger(bufio.NewWriter(os.Stdout))
}

func NewStderrLogger() Logger {
	return NewBufioLogger(bufio.NewWriter(os.Stderr))
}

func (l *BufioLogger) SetLevel(level int) {
	l.level = level
}

func (l *BufioLogger) Write(level int, message []byte) {
	if level > l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out.Write(message)
}

func (l *BufioLogger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out.Flush()
}

func (app *Application) logPrint(level int, message string) {
	if app.loggers != nil {
		name, ok := logLevelName[level]
		if ok {
			built := []byte(fmt.Sprintf("[%s] %s\n", name, message))
			for _, logger := range app.loggers {
				logger.Write(level, built)
			}
		}
	}
}

func (app *Application) logFlush() {
	if app.loggers != nil {
		for _, logger := range app.loggers {
			logger.Flush()
		}
	}
}

func (c *Context) logPrint(level int, message string) {
	c.Application.logPrint(level, message)
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

func (app *Application) SetHikaruLogLevel(level int) {
	app.hikaruLogLevel = level
}

func (app *Application) hikaruLogPrint(level int, message string) {
	if level > app.hikaruLogLevel {
		return
	}
	app.logPrint(level, message)
}

func (c *Context) hikaruLogPrint(level int, message string) {
	c.Application.hikaruLogPrint(level, message)
}

func (c *Context) logDebug(args ...interface{}) {
	c.hikaruLogPrint(LogLevelDebug, fmt.Sprint(args...))
}

func (c *Context) logDebugln(args ...interface{}) {
	c.hikaruLogPrint(LogLevelDebug, fmt.Sprintln(args...))
}

func (c *Context) logDebugf(format string, args ...interface{}) {
	c.hikaruLogPrint(LogLevelDebug, fmt.Sprintf(format, args...))
}

func (c *Context) logInfo(args ...interface{}) {
	c.hikaruLogPrint(LogLevelInfo, fmt.Sprint(args...))
}

func (c *Context) logInfoln(args ...interface{}) {
	c.hikaruLogPrint(LogLevelInfo, fmt.Sprintln(args...))
}

func (c *Context) logInfof(format string, args ...interface{}) {
	c.hikaruLogPrint(LogLevelInfo, fmt.Sprintf(format, args...))
}

func (c *Context) logWarning(args ...interface{}) {
	c.hikaruLogPrint(LogLevelWarn, fmt.Sprint(args...))
}

func (c *Context) logWarningln(args ...interface{}) {
	c.hikaruLogPrint(LogLevelWarn, fmt.Sprintln(args...))
}

func (c *Context) logWarningf(format string, args ...interface{}) {
	c.hikaruLogPrint(LogLevelWarn, fmt.Sprintf(format, args...))
}

func (c *Context) logError(args ...interface{}) {
	c.hikaruLogPrint(LogLevelError, fmt.Sprint(args...))
}

func (c *Context) logErrorln(args ...interface{}) {
	c.hikaruLogPrint(LogLevelError, fmt.Sprintln(args...))
}

func (c *Context) logErrorf(format string, args ...interface{}) {
	c.hikaruLogPrint(LogLevelError, fmt.Sprintf(format, args...))
}

func (c *Context) logCritical(args ...interface{}) {
	c.hikaruLogPrint(LogLevelCritical, fmt.Sprint(args...))
}

func (c *Context) logCriticalln(args ...interface{}) {
	c.hikaruLogPrint(LogLevelCritical, fmt.Sprintln(args...))
}

func (c *Context) logCriticalf(format string, args ...interface{}) {
	c.hikaruLogPrint(LogLevelCritical, fmt.Sprintf(format, args...))
}
