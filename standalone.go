// +build !appengine

package hikaru

import (
	"github.com/najeira/goutils/nlog"
)

var (
	applicationLogger nlog.Logger
	generalLogger     nlog.Logger
)

func SetApplicationLogger(logger nlog.Logger) {
	applicationLogger = logger
}

func SetGeneralLogger(logger nlog.Logger) {
	generalLogger = logger
}

var nlogLogLevelPrinterMap = map[int](func(nlog.Logger, string, ...interface{})){
	nlog.Critical: nlog.Logger.Criticalf,
	nlog.Error:    nlog.Logger.Errorf,
	nlog.Warn:     nlog.Logger.Warnf,
	nlog.Notice:   nlog.Logger.Noticef,
	nlog.Info:     nlog.Logger.Infof,
	nlog.Debug:    nlog.Logger.Debugf,
	nlog.Verbose:  nlog.Logger.Verbosef,
}

type envContext struct {
}

func (c *envContext) init() {
	// nothing for standalone environment
}

func (c *envContext) release() {
	// nothing for standalone environment
}

func (c *envContext) appLogf(level int, format string, args ...interface{}) {
	c.logf(applicationLogger, level, format, args...)
}

func (c *envContext) genLogf(level int, format string, args ...interface{}) {
	c.logf(generalLogger, level, format, args...)
}

func (c *envContext) logf(logger nlog.Logger, level int, format string, args ...interface{}) {
	if logger != nil {
		f, ok := nlogLogLevelPrinterMap[level]
		if ok {
			f(logger, format, args...)
		}
	}
}
