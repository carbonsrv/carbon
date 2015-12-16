-- bootscript.lua
-- This gets executed if there is no script given.

-- Small frontpage.
srv.GET("/", mw.echo(doctype()(
	tag"head"(
		tag"title" "Carbon"
	),
	tag"body"(
		tag"h1" "It works!",
		tag"br",
		"This means you launched the default carbon configuration, which serves as a simple static webserver with no autoindex or redirection."
	)
)))

-- Serve files.
srv.DefaultRoute(mw.static())
