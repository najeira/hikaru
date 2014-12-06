package hikaru

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Application struct {
	*Module
	Router         *httprouter.Router
	loggers        []Logger
	closed         chan struct{}
	hikaruLogLevel int
}

func NewApplication() *Application {
	app := &Application{
		Router:         httprouter.New(),
		loggers:        make([]Logger, 0),
		closed:         make(chan struct{}),
		hikaruLogLevel: LogLevelWarn,
	}
	app.Module = &Module{
		Handlers: nil,
		parent:   nil,
		prefix:   "/",
		app:      app,
	}
	//app.AddLogger(NewStderrLogger())
	return app
}

func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Router.ServeHTTP(w, req)
}

func (app *Application) Run(addr string) {
	// start a logger flusher
	go app.runLoggerFlusher()

	// ListenAndServe will block
	err := http.ListenAndServe(addr, app)

	// the app was closed
	close(app.closed)

	if err != nil {
		panic(err)
	}
}

func (app *Application) AddLogger(logger Logger) {
	app.loggers = append(app.loggers, logger)
}

func (app *Application) runLoggerFlusher() {
	app.hikaruLogPrint(LogLevelDebug, "start a logger flusher")
	interval := time.Second * 1
	for {
		select {
		case <-app.closed:
			// application was closed
			app.hikaruLogPrint(LogLevelDebug, "stop a logger flusher")
			break
		case <-time.After(interval):
			// flushes logs
			app.logFlush()
		}
	}
}
