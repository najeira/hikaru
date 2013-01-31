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


Route
-----

::
	app.Route("/", handleIndex)

You can use named parameter:
::
	app.Route("/blog/<id>", handleBlog)

That route will match e.g. "/blog/123" and "/blog/hello".

And Context.Val has "id" value:
::
	// c is *hikaru.Context
	id := c.Val("id")

The id will be "123" when "/blog/123", "hello" when "/blog/hello".


Result
------
Handlers should return a Result.
The result created by methods of Context like this:
::
	c.Text("Hello world.")
	c.Html("index", values)
	c.JsonText(some_obj)
	c.HtmlText(html_string)
	c.Raw(body, content_type)
	c.Redirect("/foo/bar")
	c.NotFound()
	c.Abort(503)
	c.Error(err)


Template
--------

::
	c.Html("index", values)

*/
