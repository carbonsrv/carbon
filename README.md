# carbon
Carbon is a Webserver written in [Go](https://golang.org) that uses the [Lua Scripting Language](http://www.lua.org/) for dynamic content and as a HTML template language.

# Installing
`go get -u github.com/vifino/carbon`

# Usage
## Flags
Note: Prefix all of these flags with "-". To get this list, just type `carbon -h`.

    Usage of ./carbon:
      -cert="": Certificate File for HTTPS
      -config="": Parse Config File
      -gzip=false: Use GZIP
      -host="": IP of Host to bind the Webserver on
      -http=true: Listen HTTP
      -http2=false: Enable HTTP/2
      -https=false: Listen HTTPS
      -key="": Key File for HTTPS
      -logger=true: Log Requests and Cache information
      -port=8080: Port to run Webserver on (HTTP)
      -ports=8443: Port to run Webserver on (HTTPS)
      -recovery=false: Recover from Panics
      -root=".": Path to Web Root
      -script="": Parse Lua Script as initialization
      -states=8: Number of Preinitialized Lua States
      -workers=8: Number of Worker threads.


# HTML Generation system

Let's start with an example:
```lua
doctype( -- Always start with the doctype
  tag"head"( -- Put your tags here.
    tag"title"("Hello World!")
  ),
  tag"body"( -- To put more than one tag in an existing tag, just put a comma after the inside tag and write your other tag after that.
    tag"h1"("Hello!"),
    tag"a"[{href="http://vifino.de/"}]("Link to my page!"), -- Supply arguments like that.
    link("http://vifino.de/")("Another link to my page!") -- Same thing as above, with a small helper function.
  )
)
```
This template language can be used both in static and dynamic configuration.

# Lua Script Configuration
Mainly used to configure the middleware used and to generate dynamic routes.
Example:
``` lua
srv.Use(mw.Logger())
srv.GET("/woot", mw.echo(doctype()(
        tag"body"(
            tag"h1"("woot")
        )
    ))
)
```

## Lua API:

`srv.Use(middleware)` adds a middleware to the server.

`srv.GET("path", middleware)` dispatches (GET) requests for `path` to the middleware `middleware`. (This works similar for all HTTP methods that [gin](https://github.com/gin-gonic/gin) supports.)

### Available middlewares/routes
The terms route and middleware are interchangeable here, since middlewares are basically just routes that call `context.Next()` when they are done/404., both work fine in `srv.GET` anyways.

These return middlewares that you add to the server using `srv.Use(middleware)` or `srv.GET(middleware)` for example.

___

`static.serve(prefix)` returns a static webserver. If you want to work in the directory root (`/`), use an empty string as prefix. (`""`).

`mw.new(code||function[, variables])` makes a new Lua Route: You pass a function and it gets run when a request hits. Use `content(html_content_or_template, response_code)` to send the data to the client. `variables` is a Lua Table you can use to bind variables to allow them to be accessed in the route. `{test="Testy"}` would bind the variable `test` with the content "Testy".

`mw.Lua()` will execute all Lua code at the requests path allowing you to render the template statements and generally make dynamic sites.

`mw.Logger()` is a simple Logger, logging requests and responses.

`mw.Recovery()` catches panics, you usually won't have to use it.

`mw.GZip()` enables GZip compression.

`mw.echo(content[, response_code])` echo's content as HTML. Also accepts templates, just use `mw.echo()` inplace of `content()`.

`mw.echoText(content[, response_code])` also echo's, but text only and doesn't accept the templates.

`mw.ExtRoute(table)` routes to different middlewares depending on file extension.

Example:
```lua
srv.Use(mw.ExtRoute({ -- Add the ExtRoute middleware.
        [".lua"]=mw.Lua(), -- Route every .lua file to mw.Lua
        ["***"]=static.serve("") -- Route every not routed request to static.serve
}))
```

# License
MIT
Copyright (c) 2015 Adrian Pistol
