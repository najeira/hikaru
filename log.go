package hikaru

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
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

func (c *Context) internalLogPrint(level int, message string) {
	c.Application.internalLogPrint(level, message)
}

func (c *Context) logDebug(args ...interface{}) {
	c.internalLogPrint(LogLevelDebug, fmt.Sprint(args...))
}

func (c *Context) logDebugln(args ...interface{}) {
	c.internalLogPrint(LogLevelDebug, fmt.Sprintln(args...))
}

func (c *Context) logDebugf(format string, args ...interface{}) {
	c.internalLogPrint(LogLevelDebug, fmt.Sprintf(format, args...))
}

func (c *Context) logInfo(args ...interface{}) {
	c.internalLogPrint(LogLevelInfo, fmt.Sprint(args...))
}

func (c *Context) logInfoln(args ...interface{}) {
	c.internalLogPrint(LogLevelInfo, fmt.Sprintln(args...))
}

func (c *Context) logInfof(format string, args ...interface{}) {
	c.internalLogPrint(LogLevelInfo, fmt.Sprintf(format, args...))
}

func (c *Context) logWarning(args ...interface{}) {
	c.internalLogPrint(LogLevelWarn, fmt.Sprint(args...))
}

func (c *Context) logWarningln(args ...interface{}) {
	c.internalLogPrint(LogLevelWarn, fmt.Sprintln(args...))
}

func (c *Context) logWarningf(format string, args ...interface{}) {
	c.internalLogPrint(LogLevelWarn, fmt.Sprintf(format, args...))
}

func (c *Context) logError(args ...interface{}) {
	c.internalLogPrint(LogLevelError, fmt.Sprint(args...))
}

func (c *Context) logErrorln(args ...interface{}) {
	c.internalLogPrint(LogLevelError, fmt.Sprintln(args...))
}

func (c *Context) logErrorf(format string, args ...interface{}) {
	c.internalLogPrint(LogLevelError, fmt.Sprintf(format, args...))
}

func (c *Context) logCritical(args ...interface{}) {
	c.internalLogPrint(LogLevelCritical, fmt.Sprint(args...))
}

func (c *Context) logCriticalln(args ...interface{}) {
	c.internalLogPrint(LogLevelCritical, fmt.Sprintln(args...))
}

func (c *Context) logCriticalf(format string, args ...interface{}) {
	c.internalLogPrint(LogLevelCritical, fmt.Sprintf(format, args...))
}

func (app *Application) SetLogger(logger Logger) {
	app.logger = logger
}

func (app *Application) SetInternalLogger(logger Logger) {
	app.internalLogger = logger
}

func logPrint(logger Logger, level int, message string) {
	if logger != nil {
		name, ok := logLevelName[level]
		if ok {
			built := []byte(fmt.Sprintf("[%s] %s\n", name, message))
			logger.Write(level, built)
		}
	}
}

func (app *Application) logPrint(level int, message string) {
	logPrint(app.logger, level, message)
}

func (app *Application) internalLogPrint(level int, message string) {
	logPrint(app.internalLogger, level, message)
}

func (app *Application) logFlush() {
	if app.logger != nil {
		app.logger.Flush()
	}
	if app.internalLogger != nil {
		app.internalLogger.Flush()
	}
}

func (app *Application) runLoggerFlusher(interval time.Duration) {
	app.internalLogPrint(LogLevelDebug, "start a logger flusher")
	for {
		select {
		case <-app.closed:
			// application was closed
			app.internalLogPrint(LogLevelDebug, "stop a logger flusher")
			app.logFlush()
			return
		case <-time.After(interval):
			// flushes logs
			app.logFlush()
		}
	}
}
