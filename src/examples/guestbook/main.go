package guestbook

import (
	"hikaru"
	"appengine/datastore"
	"time"
)

type Greeting struct {
	Content string    `datastore:",noindex"`
	Date    time.Time `datastore:","`
}

var app = hikaru.NewApplication()

func init() {
	app.Debug = true
	app.RouteFunc("/", handleIndex)
	app.RouteFunc("/sign", handleSign)
	app.Start()
}

func handleIndex(c hikaru.Context) hikaru.Result {
	q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
	greetings := make([]Greeting, 0, 10)
	if _, err := q.GetAll(c, &greetings); err != nil {
		return c.Abort(err)
	}
	return c.Html("index", greetings)
}

func handleSign(c hikaru.Context) hikaru.Result {
	g := Greeting{
		Content: c.Form("content"),
		Date:    time.Now(),
	}
	key := datastore.NewIncompleteKey(c, "Greeting", nil)
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		return c.Abort(err)
	}
	return c.Redirect("/")
}
