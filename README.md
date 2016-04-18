# carbon [![Travis CI](https://travis-ci.org/carbonsrv/carbon.svg?branch=master)](https://travis-ci.org/carbonsrv/carbon) [![Circle CI](https://circleci.com/gh/carbonsrv/carbon.svg?style=shield)](https://circleci.com/gh/carbonsrv/carbon) [![Go Report Card](https://goreportcard.com/badge/github.com/carbonsrv/carbon)](https://goreportcard.com/report/github.com/carbonsrv/carbon)
[![Docker Pulls](https://img.shields.io/docker/pulls/carbonsrv/carbon.svg)](https://hub.docker.com/r/carbonsrv/carbon/) [![Docker Stars](https://img.shields.io/docker/stars/carbonsrv/carbon.svg)](https://hub.docker.com/r/carbonsrv/carbon/) [![](https://badge.imagelayers.io/carbonsrv/carbon:latest.svg)](https://imagelayers.io/?images=carbonsrv/carbon:latest 'Get your own badge on imagelayers.io')

Carbon is a Lua Application Toolkit with focus on Web Servers, written in [Go](https://golang.org).

# Installing
`go get -u github.com/carbonsrv/carbon`

## Dependencies
- [LuaJIT 2.X](http://luajit.org)
- [PhysFS](https://icculus.org/physfs/)
- [go-bindata](https://github.com/jteeuwen/go-bindata)

The development headers of LuaJIT and PhysFS are also required.

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
