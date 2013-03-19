package hikaru

import (
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

type Application struct {
	Routes          []Route
	StaticDir       string // static file dir, default is "static"
	TemplateDir     string // template file dir, default is "templates"
	TemplateExt     string // template file ext, default is "html"
	Debug           bool
	LogLevel        int
	Renderers       map[string]Renderer
	Mutex           sync.RWMutex
	RequestFilters  []RequestFilter
	HandlerFilters  []HandlerFilter
	ErrorFilters    []ErrorFilter
	ResponseFilters []ResponseFilter
}

type RequestFilter func(*http.Request) Result
type HandlerFilter func(Context, Handler) Result
type ErrorFilter func(Context, Result) Result
type ResponseFilter func(Context, Result) Result

type Filter interface {
	Execute(...interface{}) Result
}

func NewApplication() *Application {
	app := new(Application)
	app.StaticDir = "static"
	app.TemplateDir = "templates"
	app.TemplateExt = "html"
	app.LogLevel = LogLevelInfo
	app.Renderers = make(map[string]Renderer)
	app.RequestFilters = make([]RequestFilter, 0)
	app.HandlerFilters = make([]HandlerFilter, 0)
	app.ErrorFilters = make([]ErrorFilter, 0)
	app.ResponseFilters = make([]ResponseFilter, 0)
	return app
}

func (app *Application) Start() {
	app.initRenderers()
	http.Handle("/", app)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	done := app.handlePprof(w, r)
	if done {
		return
	}
	c := NewContext(app, w, r)
	c.Execute()
}

func (app *Application) GetRenderer(kind string) Renderer {
	app.Mutex.RLock()
	defer app.Mutex.RUnlock()
	r, _ := app.Renderers[kind]
	return r
}

func (app *Application) SetRenderer(kind string, r Renderer) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	app.Renderers[kind] = r
}

func (app *Application) AddRequestFilter(h RequestFilter) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	app.RequestFilters = append(app.RequestFilters, h)
}

func (app *Application) AddHandlerFilter(h HandlerFilter) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	app.HandlerFilters = append(app.HandlerFilters, h)
}

func (app *Application) AddErrorFilter(h ErrorFilter) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	app.ErrorFilters = append(app.ErrorFilters, h)
}

func (app *Application) AddResponseFilter(h ResponseFilter) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	app.ResponseFilters = append(app.ResponseFilters, h)
}

func (app *Application) initRenderers() {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	if r, _ := app.Renderers["html"]; r == nil {
		app.Renderers["html"] = NewRenderer(app.TemplateDir, app.TemplateExt)
	}
}

func (app *Application) RouteFunc(pattern string, handler Handler) {
	route := NewRoute(pattern, handler)
	app.Route(route)
}

func (app *Application) RouteFuncTimeout(pattern string, handler Handler, timeout time.Duration) {
	route := NewRoute(pattern, handler)
	route.timeout = timeout
	app.Route(route)
}

func (app *Application) Route(route Route) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()
	app.Routes = append(app.Routes, route)
}

func (app *Application) Match(r *http.Request) *RouteData {
	app.Mutex.RLock()
	defer app.Mutex.RUnlock()
	var rd *RouteData
	for _, route := range app.Routes {
		rd = route.Match(r)
		if rd != nil {
			return rd
		}
	}
	return nil
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
