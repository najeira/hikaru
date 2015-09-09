package hikaru

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path"
)

type HandlerFunc func(*Context)

type Module struct {
	handlers []HandlerFunc
	prefix   string
	parent   *Module
	app      *Application
	router   *httprouter.Router
}

func (m *Module) Use(middlewares ...HandlerFunc) {
	m.handlers = append(m.handlers, middlewares...)
}

func (m *Module) Module(component string, handlers ...HandlerFunc) *Module {
	prefix := path.Join(m.prefix, component)
	return &Module{
		handlers: m.combineHandlers(handlers),
		parent:   m,
		prefix:   prefix,
		router:   m.router,
	}
}

func (m *Module) Handle(method, p string, handlers []HandlerFunc) {
	var h handlerFuncs = m.combineHandlers(handlers)
	m.router.Handle(method, path.Join(m.prefix, p), h.handle)
}

type handlerFuncs []HandlerFunc

func (h handlerFuncs) handle(w http.ResponseWriter, r *http.Request, hp httprouter.Params) {
	c := getContext()
	defer releaseContext(c)

	c.init(w, r, h)
	c.addParams(hp)
	//c.verbosef("hikaru: Request is %v", c.Request)
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
	s := len(m.handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	if m.handlers != nil && len(m.handlers) > 0 {
		h = append(h, m.handlers...)
	}
	if handlers != nil && len(handlers) > 0 {
		h = append(h, handlers...)
	}
	return h
}
