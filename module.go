package hikaru

import (
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

type HandlerFunc func(*Context)

type Module struct {
	parent *Module
	prefix string
	router *httprouter.Router
}

func NewModule(prefix string) *Module {
	return &Module{
		parent: nil,
		prefix: prefix,
		router: httprouter.New(),
	}
}

func (m *Module) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.router.ServeHTTP(w, req)
}

func (m *Module) Module(component string, handler HandlerFunc) *Module {
	prefix := path.Join(m.prefix, component)
	return &Module{
		parent: m,
		prefix: prefix,
		router: m.router,
	}
}

func (m *Module) Handle(method, p string, handler HandlerFunc) {
	h := wrapHandlerFunc(handler)
	m.router.Handle(method, path.Join(m.prefix, p), h)
}

func wrapHandlerFunc(handler HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, hp httprouter.Params) {
		c := getContext(w, r, hp)
		defer releaseContext(c)
		//c.verbosef("hikaru: Request is %v", c.Request)
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
	p = path.Join(p, "/*filepath")
	fileServer := http.FileServer(http.Dir(root))
	m.GET(p, func(c *Context) {
		fp, err := c.TryString("filepath")
		if err != nil {
			c.Response.Fail(err)
		} else {
			original := c.Request.URL.Path
			c.Request.URL.Path = fp
			fileServer.ServeHTTP(c.Response, c.Request)
			c.Request.URL.Path = original
		}
	})
}
