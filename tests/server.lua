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
	content(type("wee"), 200)
end))
srv.POST("/", mw.new(function()
	print(form("f"))
	content(form("f"))
	if form("f") == nil then
		print("is nil")
	end
end))
srv.GET("/bindtest", mw.new(function()
	test(context)
end, {test=mw.echo("test")}))
srv.GET("/ws", mw.ws(function()
	while true do
		print(ws.read())
	end
end))
print("bai")
