package hikaru

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Application struct {
	*Module
}

func New() *Application {
	app := &Application{
		Module: &Module{
			handlers: nil,
			parent:   nil,
			prefix:   "/",
			router:   httprouter.New(),
		},
	}
	return app
}

func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Module.router.ServeHTTP(w, req)
}

func (app *Application) Run(addr string) {
	// ListenAndServe will block
	err := http.ListenAndServe(addr, app)

	if err != nil {
		panic(err)
	}
}
