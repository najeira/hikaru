package hikaru

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Application struct {
	*Module
	Router         *httprouter.Router
	logger         Logger
	internalLogger Logger
	closed         chan struct{}
}

func New() *Application {
	logger := NewStderrLogger()
	logger.SetLevel(LogLevelWarn)
	app := &Application{
		Router:         httprouter.New(),
		closed:         make(chan struct{}),
		internalLogger: logger,
	}
	app.Module = &Module{
		Handlers: nil,
		parent:   nil,
		prefix:   "/",
		app:      app,
	}
	return app
}

func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Router.ServeHTTP(w, req)
}

func (app *Application) Run(addr string) {
	// start a logger flusher
	go app.runLoggerFlusher(time.Second * 1)

	// ListenAndServe will block
	err := http.ListenAndServe(addr, app)

	// the app was closed
	close(app.closed)

	if err != nil {
		panic(err)
	}
}
