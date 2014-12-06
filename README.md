# About

Hikaru is a web framework for Go.
It supports standalone and Google App Engine.

*This is under construction.*
Do not use for production services.


# Getting started

Hello World:

```
package main

import "github.com/najeira/hikaru"

func main() {
    app := hikaru.New()
    app.GET("/", func(c *hikaru.Context) {
        c.Text("Hello World")
    })
    app.Run(":8080")
}
```


# Routing

Using GET, POST, HEAD and OPTIONS:

```
app.GET("/get", get)
app.POST("/post", post)
app.HEAD("/head", head)
app.OPTIONS("/options", options)
```

Hikaru uses julienschmidt's httprouter internaly.

See: https://github.com/julienschmidt/httprouter

# Request

## Parameters

Query parameters:

```
app.GET("/user", func(c *hikaru.Context)) {
    // String get "Gopher" when the query is /user?name=Gopher
    // Second argument of String() is failover that will return 
    // when the name does not exists in the query
    name := c.String("name", "")
    c.Text("Hello " + name)
})
```

Getting parameters in path same as query:

```
app.GET("/user/:name", func(c *hikaru.Context)) {
    // String get "Gopher" when the path is /user/Gopher
    name := c.String("name", "")
    c.Text("Hello " + name)
})
```

Shorthand to get int and float parameters:

```
app.GET("/user/:id", func(c *hikaru.Context)) {
    // Returns 123 if /user/123
    id := c.Int("id", 0) // id is int64
})
app.GET("/hoge/:score", func(c *hikaru.Context)) {
    // Returns 12.3 if /hoge/12.3
    score := c.Float("score", 0.0) // score is float64
})
```

# Response

## Writing response

```
// Text writes string
c.Text("Hello world.")

// Json marshals the object to json string and write
c.Json(someObj)

// Raw writes []byte
c.Raw(body, "application/octet-stream")
```

## Redirect

```
// RedirectFound sends 302 Found
c.RedirectFound("/foo/bar")

// RedirectFound sends 301 Moved Permanently
c.RedirectMoved("/foo/bar")

// Redirect sends location and code
c.Redirect("/foo/bar", statusCode)
```

## Error response

```
// 304 Not Modified
c.NotModified()

// 401 Unauthorized
c.Unauthorized()

// 403 Forbidden
c.Forbidden()

// HTTP 404 Not Found
c.NotFound()

// 500 Internal Server Error
c.Fail(err)
```

You can send your body with error status:

```
c.NotFound()
c.Json(errObj)
```

## Headers

```
// Header gets the response headers.
headers := c.Header()

// SetHeader sets a header value to the response.
c.SetHeader("Foo", "Bar")

// SetHeader sets a header value to the response.
c.SetHeader("Foo", "Bar")
```
