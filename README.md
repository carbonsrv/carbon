# gold
Go webserver using Lua for dynamic content.

# About
This is a go Webserver that is using Lua for dynamic content, either with no configuration in a plain 
serve-pwd-config, running every .lua that gets accessed, configured using a lua script or just command line flags.

# Usage

## Flags
    Usage of ./gold:
      -config="": Parse Config File
      -gzip=false: Use GZIP
      -host="": IP of Host to bind the Webserver on
      -logger=true: Log Requests and Cache information
      -port=8080: Port to run Webserver on (HTTP)
      -recovery=false: Recover from Panics
      -root=".": Path to Web Root
      -script="": Parse Lua Script as initialization
      -states=8: Number of Preinitialized Lua States
      -workers=8: Number of Worker threads.

## Script configuration
There is also the script configuration, which uses a lua script to do the initialization.
More on that later.

# HTML Generation system
As you could spy in that script config example, there is a very interesting alternative to writing plain html.

`content(cont[, response_code])` sets the content in the http request, but it takes more than just strings, it also takes a tag object.

To start with the basics:
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
Where is that thing used?
Well, either in the Lua files run by a non-script configured server, or in the Lua Dynamic Routes and Lua Script configuration file.

# Lua Script Configuration
Mainly used to configure the middleware used and to generate dynamic routes.
Example:
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
Now in a little more detail.

`srv.Use(middlware)` adds a middleware to the server.

### Available middlewares
`mw.Lua()` behaves a little like PHP, it runs every file it gets, using the same `content(...)` you can make dynamic sites!

`mw.Logger()` is a fancy logger. 'Nuff said.

`mw.Recovery()` is an panic catcher, normally not used.

`mw.GZip()` GZips everything that goes through it.

`mw.new(code||function)` makes a new lua route. Look above for an example on how to use it.

`mw.Echo(response_code, content)` echo's a string. Simple as that.

`mw.EchoHTML(response_code, content)` is the same as above, but don't uses plaintext and uses HTML instead.

`static.serve(prefix)` is the most basic webserver. Static. Use `""` as prefix for `/`, please.

`mw.ExtRoute(table)` is one of the more interesting middlewares. It routes to different middlewares depending on file extension.
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
