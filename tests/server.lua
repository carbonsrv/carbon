print("hai")
srv.Use(mw.Logger())
srv.Use(mw.ExtRoute({
	[".lua"]=mw.Lua(),
	["***"]=static.serve("")
}))
srv.GET("/woot", mw.echo(doctype()(
		tag"body"(
			tag"h1"("woot")
		)
), 200))
srv.GET("/wat", mw.echoText("u wut m8?!?!"))
srv.GET("/test", mw.new(function()
	content(test, 200)
end, 200, {test="wat"}))
print("bai")
