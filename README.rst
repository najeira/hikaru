About
=====

Hikaru is a web framework for Google App Engine Go.

This is under construction. Do not use production.


Getting started
===============

Hello World:
::
    package hello
    
    import "github.com/najeira/hikaru"
    
    var app = hikaru.NewApplication()
    
    func init() {
    	app.RouteFunc("/", handleWelcome)
    	app.Start()
    }
    
    func handleWelcome(c hikaru.Context) hikaru.Result {
    	return c.Text("Google App Engine Go!")
    }


Route
=====

Bind func to path:
::
    app.RouteFunc("/", handleWelcome)


Named parameter
---------------

To use named parameter, use '<' and '>':
::
    app.RouteFunc("/blog/<id>", handleBlog)

That route will match e.g. "/blog/123" and "/blog/hello".

And Context.Val has "id" value:
::
    id := c.Val("id")

The id will be "123" when "/blog/123", "hello" when "/blog/hello".


Custom route
------------

To use your original route, implements hikaru.Route interface and it pass to Route method:
::
    app.Route(route)


Result
======

Handlers should return a Result.
The result created by methods of Context like this:
::
    c.Text("Hello world.")
    c.Html("index", values)
    //c.JsonText(some_obj)
    //c.HtmlText(html_string)
    c.Raw(body, content_type)
    c.Redirect("/foo/bar")
    c.NotFound()
    c.Abort(err)
    c.AbortCode(503)


Template
========

To Render html template:
::
    c.Html("index", values)
