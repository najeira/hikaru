package hikaru

type Application struct {
	*Module
}

func New() *Application {
	return &Application{
		Module: NewModule("/"),
	}
}
