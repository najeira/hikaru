package hikaru

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Config struct {
	ProxyAddr string
}

type Application struct {
	*Module
	Config       *Config
	Router       *httprouter.Router
	loggers      []Logger
	hikaruLogger Logger
}

func New(config *Config) *Application {
	if config == nil {
		config = &Config{}
	}
	app := &Application{
		Config:  config,
		Router:  httprouter.New(),
		loggers: make([]Logger, 0),
	}
	app.SetHikaruLog(LogLevelWarn)
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
	if err := http.ListenAndServe(addr, app); err != nil {
		panic(err)
	}
}

func (app *Application) AddLogger(logger Logger) {
	app.loggers = append(app.loggers, logger)
}
