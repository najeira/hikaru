// +build appengine

package hikaru

import (
	"appengine"
	"github.com/najeira/goutils/nlog"
)

var (
	applicationLogLevel int
	generalLogLevel     int
)

func SetApplicationLogLevel(level int) {
	applicationLogLevel = level
}

func SetGeneralLogLevel(level int) {
	generalLogLevel = level
}

var appEngineLogLevelPrinterMap = map[int](func(appengine.Context, int, string, ...interface{})){
	nlog.Critical: appengine.Context.Criticalf,
	nlog.Error:    appengine.Context.Errorf,
	nlog.Warn:     appengine.Context.Warningf,
	nlog.Notice:   appengine.Context.Infof,
	nlog.Info:     appengine.Context.Infof,
	nlog.Debug:    appengine.Context.Debugf,
	nlog.Verbose:  appengine.Context.Debugf,
}

type envContext struct {
	appengine.Context
}

func (c *envContext) init() {
	c.Context = appengine.NewContext(c.Request)
}

func (c *envContext) release() {
	c.Context = nil
}

func (c *envContext) appLogf(level int, format string, args ...interface{}) {
	if level > applicationLogLevel {
		c.logf(level, format, args...)
	}
}

func (c *envContext) genLogf(level int, format string, args ...interface{}) {
	if level > generalLogLevel {
		c.logf(level, format, args...)
	}
}

func (c *envContext) logf(level int, format string, args ...interface{}) {
	f, ok := appEngineLogLevelPrinterMap[level]
	if ok {
		f(c.Context, format, args...)
	}
}
