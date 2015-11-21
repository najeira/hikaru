package hikaru

import (
	"net/http"
)

type Application struct {
	*Module
}

func New() *Application {
	return &Application{
		Module: NewModule("/"),
	}
}

func (app *Application) Run(addr string) {
	// ListenAndServe will block
	err := http.ListenAndServe(addr, app)

	if err != nil {
		panic(err)
	}
}
