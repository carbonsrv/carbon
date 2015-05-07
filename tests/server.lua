print("hai")
srv.Use(mw.Logger())
srv.Use(mw.ExtRoute({
	[".lua"]=mw.Lua(),
	["***"]=static.serve("")
}))
srv.GET("/woot", mw.new(function()
	content(doctype()(
		tag"body"(
			tag"h1"("woot")
		)
	))
end))
print("bai")
