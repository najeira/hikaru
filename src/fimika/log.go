package fimika

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	LogLevelNo = iota
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

type Logger struct {
	LogLevel int
	logger   *log.Logger
}

func NewLogger() *Logger {
	lg := log.New(os.Stdout, "", log.LstdFlags)
	df := &Logger{logger: lg, LogLevel: LogLevelWarn}
	return df
}

func (l *Logger) SetOutput(out io.Writer) {
	l.logger = log.New(out, l.logger.Prefix(), l.logger.Flags())
}

func (l *Logger) Flags() int {
	return l.logger.Flags()
}

func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

func (l *Logger) Prefix() string {
	return l.logger.Prefix()
}

func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *Logger) Print(level int, args ...interface{}) {
	if l.LogLevel >= level {
		l.logger.Output(3, fmt.Sprint(args...))
	}
}

func (l *Logger) Println(level int, args ...interface{}) {
	if l.LogLevel >= level {
		l.logger.Output(3, fmt.Sprintln(args...))
	}
}

func (l *Logger) Printf(level int, format string, args ...interface{}) {
	if l.LogLevel >= level {
		l.logger.Output(3, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.Print(LogLevelDebug, fmt.Sprint(args...))
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Println(LogLevelDebug, fmt.Sprintln(args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Printf(LogLevelDebug, fmt.Sprintf(format, args...))
}

func buildMessage(prefix string, args ...interface{}) []interface{} {
	v := make([]interface{}, 1, len(args) + 1)
	v[0] = prefix
	v = append(v, args...)
	return v
}

func (l *Logger) Info(args ...interface{}) {
	l.Print(LogLevelInfo, buildMessage("[INFO]", args...)...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Println(LogLevelInfo, buildMessage("[INFO]", args...)...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Printf(LogLevelInfo, "[INFO] " + format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Print(LogLevelWarn, buildMessage("[WARN]", args...)...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.Println(LogLevelWarn, buildMessage("[WARN]", args...)...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Printf(LogLevelWarn, "[WARN] " + format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Print(LogLevelError, buildMessage("[ERROR]", args...)...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Println(LogLevelError, buildMessage("[ERROR]", args...)...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Printf(LogLevelError, "[ERROR] " + format, args...)
}
