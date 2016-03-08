package hikaru

import (
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

type HandlerFunc func(*Context)

type Middleware func(HandlerFunc) HandlerFunc

type Module struct {
	parent     *Module
	prefix     string
	router     *httprouter.Router
	middleware Middleware
}

func NewModule(prefix string) *Module {
	return &Module{
		parent:     nil,
		prefix:     prefix,
		router:     httprouter.New(),
		middleware: nil,
	}
}

func (m *Module) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.router.ServeHTTP(w, req)
}

func (m *Module) Module(component string, middleware Middleware) *Module {
	prefix := path.Join(m.prefix, component)
	return &Module{
		parent:     m,
		prefix:     prefix,
		router:     m.router,
		middleware: middleware,
	}
}

func (m *Module) Handle(method, p string, handler HandlerFunc) {
	if m.middleware != nil {
		handler = m.middleware(handler)
	}
	h := m.wrapHandlerFunc(handler)
	m.router.Handle(method, path.Join(m.prefix, p), h)
}

func (m *Module) wrapHandlerFunc(handler HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, hp httprouter.Params) {
		c := getContext(w, r, hp)
		defer releaseContext(c)
		c.tracef("hikaru: Request is %v", c.Request)
		handler(c)
	}
}

// POST is a shortcut for Module.Handle("POST", p, handle)
func (m *Module) POST(p string, handler HandlerFunc) {
	m.Handle("POST", p, handler)
}

// GET is a shortcut for Module.Handle("GET", p, handle)
func (m *Module) GET(p string, handler HandlerFunc) {
	m.Handle("GET", p, handler)
}

// OPTIONS is a shortcut for Module.Handle("OPTIONS", p, handle)
func (m *Module) OPTIONS(p string, handler HandlerFunc) {
	m.Handle("OPTIONS", p, handler)
}

// HEAD is a shortcut for Module.Handle("HEAD", p, handle)
func (m *Module) HEAD(p string, handler HandlerFunc) {
	m.Handle("HEAD", p, handler)
}

func (m *Module) Static(p, root string) {
	p = path.Join(p, "*filepath")
	fileServer := http.FileServer(http.Dir(root))
	m.GET(p, func(c *Context) {
		fp, err := c.TryString("filepath")
		if err != nil {
			c.Errorf(err.Error())
			c.Fail()
		} else {
			original := c.Request.URL.Path
			c.Request.URL.Path = fp
			fileServer.ServeHTTP(c.ResponseWriter, c.Request)
			c.Request.URL.Path = original
		}
	})
}
