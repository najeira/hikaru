package hikaru

import (
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

type Application struct {
	Routes         []*Route
	StaticDir      string // static file dir, default is "static"
	TemplateDir    string // template file dir, default is "templates"
	TemplateExt    string // template file ext, default is "html"
	HandlerTimeout time.Duration
	Debug          bool
	LogLevel       int
	Renderer       Renderer
	Mutex          sync.RWMutex
}

func NewApplication() *Application {
	app := new(Application)
	app.Routes = make([]*Route, 0)
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
	done := app.handlePprof(w, r)
	if done {
		return
	}
	c := NewContext(app, w, r)
	ok := c.executeRoute()
	if ok {
		c.executeContext()
	} else {
		c.executeNotFound()
	}
	c.executeResult()
}

func (app *Application) SetRenderer(r Renderer) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	app.Renderer = r
}

func (app *Application) InitRenderer() {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	if app.Renderer == nil {
		app.Renderer = NewRenderer(app.TemplateDir, app.TemplateExt)
	}
}

func (app *Application) Route(pattern string, handler Handler) {
	route := NewRoute(pattern, handler)
	app.appendRoute(route)
}

func (app *Application) RouteTimeout(pattern string, handler Handler, timeout time.Duration) {
	route := NewRoute(pattern, handler)
	route.Timeout = timeout
	app.appendRoute(route)
}

func (app *Application) appendRoute(route *Route) {
	app.Routes = append(app.Routes, route)
}

func (app *Application) handlePprof(w http.ResponseWriter, r *http.Request) bool {
	if !app.Debug {
		return false
	}
	switch r.URL.Path {
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
