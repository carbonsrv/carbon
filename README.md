# carbon [![Travis CI](https://travis-ci.org/carbonsrv/carbon.svg?branch=master)](https://travis-ci.org/carbonsrv/carbon) [![Circle CI](https://circleci.com/gh/carbonsrv/carbon.svg?style=shield)](https://circleci.com/gh/carbonsrv/carbon) [![Docker Pulls](https://img.shields.io/docker/pulls/carbonsrv/carbon.svg)](https://hub.docker.com/r/carbonsrv/carbon/) [![Docker Stars](https://img.shields.io/docker/stars/carbonsrv/carbon.svg)](https://hub.docker.com/r/carbonsrv/carbon/) 

Carbon is a Webserver written in [Go](https://golang.org) that uses the [Lua Scripting Language](http://www.lua.org/) for dynamic content and as a HTML template language.

# Installing
`go get -u github.com/carbonsrv/carbon`

# Usage
## Flags
To get this list, just type `carbon -h`.

    Usage of carbon:
      -cert="": Certificate File for HTTPS
      -config="": Parse Config File
      -debug=false: Show debug information
      -eval="": Eval Lua Code
      -gzip=false: Use GZIP
      -host="": IP of Host to bind the Webserver on
      -http=true: Listen HTTP
      -http2=false: Enable HTTP/2
      -https=false: Listen HTTPS
      -key="": Key File for HTTPS
      -logger=true: Log Request information
      -port=8080: Port to run Webserver on (HTTP)
      -ports=8443: Port to run Webserver on (HTTPS)
      -recovery=false: Recover from Panics
      -repl=false: Run REPL
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

# More info

Check the wiki!

# License
MIT
Copyright (c) 2015 Adrian Pistol

Third party software included with this may have different licenses. Check /builtin/NOTICE.txt.
