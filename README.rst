About
=====

Hikaru is a web framework for Go supports standalone and Google AppEngine.

This is under construction. Do not use for production services.


Getting started
===============

Hello World:
::
    package main
    
    import "github.com/najeira/hikaru"
    
    func main() {
        app := hikaru.New(nil)
        app.GET("/", func(c *hikaru.Context) {
            c.Text("Hello World")
        })
        app.Run(":8080")
    }


Route
=====

See: https://github.com/julienschmidt/httprouter


Result
======

Handlers should return a Result.
The result created by methods of Context like this:
::
    c.Text("Hello world.")
    c.Json(some_obj)
    c.Raw(body, content_type)
    c.Redirect("/foo/bar")
    c.NotFound()
    c.Abort(503)
