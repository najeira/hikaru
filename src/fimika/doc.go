package fimika
/*

Getting started
===============

::
	package hello
	
	import "fimika"
	
	var app = fimika.NewApplication()
	
	func init() {
		app.Route("/", handleIndex)
		app.Start()
	}
	
	func handleIndex(c *fimika.Context) *fimika.Result {
		return c.Text("Hello Fimika!")
	}


Result
------
Handlers should return a Result.
The result id created by Context's methods like this:
::
	c.Text("Hello world.")
	c.JsonText(some_obj)
	c.HtmlText(html_string)
	c.Raw(binary, content_type)
	c.Redirect("/foo/bar")
	c.Error(err)


Template
--------

::
	c.Html("index", values)

*/
