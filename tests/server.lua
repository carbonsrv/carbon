print("hai")
Use(mw.Logger())
Use(mw.ExtRoute({
	[".lua"]=mw.Lua(),
	["***"]=static.serve("")
}))
print("bai")
