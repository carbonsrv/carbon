print("hai")
srv.Use(mw.Logger())
srv.Use(mw.ExtRoute({
	[".lua"]=mw.Lua(),
	["***"]=static.serve("")
}))
srv.GET("/woot", mw.new([[
context.String(200,"woot")
]]))
print("bai")
