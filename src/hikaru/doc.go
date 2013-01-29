package hikaru
/*

Getting started
===============

::
	package hello
	
	import "hikaru"
	
	var app = hikaru.NewApplication()
	
	func init() {
		app.Route("/", handleIndex)
		app.Start()
	}
	
	func handleIndex(c *hikaru.Context) hikaru.Result {
		return c.Text("Hello Hikaru!")
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
