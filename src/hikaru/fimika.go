package fimika

import (
	"appengine"
	"bytes"
	"fmt"
	"net/http"
	"net/http/pprof"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Application struct {
	Routes         []*Route
	RootDir        string // project root dir
	StaticDir      string // static file dir, "static" if empty
	TemplateDir    string // template file dir, "templates" if empty
	TemplateExt    string // template file ext, "html" if empty
	HandlerTimeout time.Duration
	Debug          bool
	LogLevel       int
	Renderer       Renderer
	mutex          sync.RWMutex
}

func NewApplication() *Application {
	app := new(Application)
	app.Routes = make([]*Route)
	app.StaticDir = "static"
	app.TemplateDir = "templates"
	app.TemplateExt = "html"
	app.Debug = false
	app.LogLevel = LogLevelInfo
	return app
}

func (app *Application) Start() {
	app.InitRenderer()
	http.Handle("/", app)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(app, w, r)

	c.RouteData = app.matchRoute(r)

	if c.RouteData == nil {
		c.Result = c.NotFound()
	} else {
		app.executeContext(c)
	}

	c.executeResult()
}

func (app *Application) matchRoute(r *http.Request) *RouteData {
	var rd *RouteData
	for _, route := range app.Routes {
		rd = route.Match(r)
		if rd != nil {
			return rd
		}
	}
	return nil
}

func (app *Application) executeContext(c *Context) {
	defer func() {
		err_result := app.recoverPanic(c)
		if err_result != nil {
			c.Result = err_result
		}
	}()

	//TODO: before handler middlewares

	c.executeHandler()

	//TODO: after handler middlewares
}

func (app *Application) recoverPanic(c *Context) Resulter {
	err := recover()
	if err == nil {
		return nil
	}

	var buf bytes.Buffer
	buf.Write(debug.Stack())
	stack := buf.String()

	err_msg := fmt.Sprintf("%v\n%s", err, stack)

	c.LogErrorf(err_msg)

	result := c.Error(err)
	if app.Debug {
		result.Body.WriteString(err_msg)
	}
	return result
}

func (app *Application) SetRenderer(r Renderer) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.Renderer = r
}

func (app *Application) InitRenderer() {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.Renderer == nil {
		app.Renderer = NewRenderer(app.TemplateDir, app.TemplateExt)
	}
}

func (app *Application) Route(pattern string, handler Handler) {
	route := NewRoute(pattern, handler)
	app.Routes = append(app.Routes, route)
}

func (app *Application) RouteTimeout(pattern string, handler Handler, timeout time.Duration) {
	route := NewRoute(pattern, handler)
	route.Timeout = timeout
	app.Routes = append(app.Routes, route)
}

func (app *Application) handlePprof(c *Context) bool {
	w := c.ResponseWriter
	r := c.Request
	if !app.Debug {
		return false
	}
	switch r.RequestURI {
	case "/debug/pprof/cmdline":
		pprof.Cmdline(w, r)
		return true
	case "/debug/pprof/profile":
		pprof.Profile(w, r)
		return true
	case "/debug/pprof/heap":
		h := pprof.Handler("heap")
		h.ServeHTTP(w, r)
		return true
	case "/debug/pprof/symbol":
		pprof.Symbol(w, r)
		return true
	default:
		return false
	}
	return false
}
