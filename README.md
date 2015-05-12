# gold
gold is a Webserver written in [Go](https://golang.org) that uses the [Lua Scripting Language](http://www.lua.org/) for dynamic content and as a HTML template language.

# Usage

## Flags
    Usage of ./gold:
      -cert="": Certificate File for HTTPS
      -config="": Parse Config File
      -gzip=false: Use GZIP
      -host="": IP of Host to bind the Webserver on
      -http=true: Listen HTTP
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
Example (Same as above):
``` lua
srv.Use(mw.Logger())
srv.GET("/woot", mw.new(function()
        content(doctype()(
                tag"body"(
                        tag"h1"("woot")
                )
        ))
end))
```

## Lua API:

`static.serve(prefix)` Starts a static webserver. If you want to work in the directory root (`/`), use an empty string as prefix. (`""`).

`srv.Use(middleware)` adds a middleware to the server.

`srv.GET("path", middleware)` dispatches (GET) requests for `path` to the middleware `middleware`. (This works similar for all HTTP methods that [gin](https://github.com/gin-gonic/gin) supports.)

### Available middlewares

`mw.new(code||function)` makes a new lua route. Look above for an example on how to use it.

`mw.Lua()` will execute all Lua code in a given file and renders the template statements.

`mw.Logger()` is a simple Logger, logging requests and responses.

`mw.Recovery()` catches panics, you usually won't have to use it.

`mw.GZip()` enables HTTP compression.

`mw.Echo(response_code, content)` echo's content as string.

`mw.EchoHTML(response_code, content)` also echo's, but renders the content to HTML first.

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
