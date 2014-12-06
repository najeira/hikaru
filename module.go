package hikaru

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path"
)

type HandlerFunc func(*Context)

type Module struct {
	Handlers []HandlerFunc
	prefix   string
	parent   *Module
	app      *Application
}

func (m *Module) Use(middlewares ...HandlerFunc) {
	m.Handlers = append(m.Handlers, middlewares...)
}

func (m *Module) Module(component string, handlers ...HandlerFunc) *Module {
	prefix := path.Join(m.prefix, component)
	return &Module{
		Handlers: m.combineHandlers(handlers),
		parent:   m,
		prefix:   prefix,
		app:      m.app,
	}
}

func (m *Module) Handle(method, p string, h []HandlerFunc) {
	pp := path.Join(m.prefix, p)
	h = m.combineHandlers(h)
	m.app.Router.Handle(method, pp, func(w http.ResponseWriter, r *http.Request, hp httprouter.Params) {
		m.Execute(w, r, h, hp)
	})
}

func (m *Module) Execute(w http.ResponseWriter, r *http.Request, h []HandlerFunc, hp httprouter.Params) {
	c := getContext(m.app, w, r, h)
	defer releaseContext(c)
	if hp != nil {
		for _, v := range hp {
			c.Add(v.Key, v.Value)
		}
	}
	c.logDebugf("execute: url is %v", c.Request.URL)
	c.execute()
}

// POST is a shortcut for Module.Handle("POST", p, handle)
func (m *Module) POST(p string, handlers ...HandlerFunc) {
	m.Handle("POST", p, handlers)
}

// GET is a shortcut for Module.Handle("GET", p, handle)
func (m *Module) GET(p string, handlers ...HandlerFunc) {
	m.Handle("GET", p, handlers)
}

// OPTIONS is a shortcut for Module.Handle("OPTIONS", p, handle)
func (m *Module) OPTIONS(p string, handlers ...HandlerFunc) {
	m.Handle("OPTIONS", p, handlers)
}

// HEAD is a shortcut for Module.Handle("HEAD", p, handle)
func (m *Module) HEAD(p string, handlers ...HandlerFunc) {
	m.Handle("HEAD", p, handlers)
}

func (m *Module) Static(p, root string) {
	p = path.Join(p, "/*filepath")
	fileServer := http.FileServer(http.Dir(root))
	m.GET(p, func(c *Context) {
		fp, err := c.TryString("filepath")
		if err != nil {
			c.Fail(err)
		} else {
			original := c.Request.URL.Path
			c.Request.URL.Path = fp
			fileServer.ServeHTTP(c.ResponseWriter, c.Request)
			c.Request.URL.Path = original
		}
	})
}

func (m *Module) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(m.Handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	if m.Handlers != nil && len(m.Handlers) > 0 {
		h = append(h, m.Handlers...)
	}
	if handlers != nil && len(handlers) > 0 {
		h = append(h, handlers...)
	}
	return h
}
