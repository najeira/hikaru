package fimika

import (
	"appengine"
	"strings"
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

type Logger struct {
	LogLevel         int
	AppEngineContext appengine.Context
}

func NewLogger(level int, c appengine.Context) *Logger {
	return &Logger{LogLevel: level, AppEngineContext: c}
}

func (l *Logger) Print(level int, args ...interface{}) {
	format := strings.Repeat("%v", len(args))
	l.Printf(level, format, args...)
}

func (l *Logger) Println(level int, args ...interface{}) {
	format := strings.Repeat("%v", len(args))
	format += "\n"
	l.Printf(level, format, args...)
}

func (l *Logger) Printf(level int, format string, args ...interface{}) {
	if l.LogLevel >= level {
		f, ok := levelFuncMap[level]
		if ok && f != nil {
			f(l.AppEngineContext, format, args...)
		}
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.Print(LogLevelDebug, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Println(LogLevelDebug, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Printf(LogLevelDebug, format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Print(LogLevelInfo, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Println(LogLevelInfo, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Printf(LogLevelInfo, format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Print(LogLevelWarn, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.Println(LogLevelWarn, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Printf(LogLevelWarn, format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Print(LogLevelError, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Println(LogLevelError, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Printf(LogLevelError, format, args...)
}
