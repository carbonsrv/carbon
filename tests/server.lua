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
	content(type(new), 200)
end))
srv.POST("/", mw.new(function()
	print(form("f"))
	content(form("f"))
	if form("f") == nil then
		print("is nil")
	end
end))
print("bai")
