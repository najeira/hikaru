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
	h = m.combineHandlers(h)
	m.app.Router.Handle(method, p, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		m.Execute(w, r, h, p)
	})
}

func (m *Module) Execute(w http.ResponseWriter, r *http.Request, h []HandlerFunc, p httprouter.Params) {
	c := NewContext(m.app, w, r, h)
	if p != nil {
		for _, v := range p {
			c.Values.Add(v.Key, v.Value)
		}
	}
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
		fp, err := c.Values.String("filepath")
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
	h = append(h, m.Handlers...)
	h = append(h, handlers...)
	return h
}
