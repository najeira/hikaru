package hikaru

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	Debug     bool
	LogLevel  int
	ProxyAddr string
	Timeout   time.Duration
}

type Application struct {
	*Module
	Config *Config
	Router *httprouter.Router
	mutex  sync.RWMutex
}

func New(config *Config) *Application {
	if config == nil {
		config = &Config{LogLevel: LogLevelInfo}
	}
	app := &Application{
		Config: config,
		Router: httprouter.New(),
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
	if err := http.ListenAndServe(addr, app); err != nil {
		panic(err)
	}
}
