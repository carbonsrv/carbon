-- Init Glue for Config State
kvstore.set("tmp:msgpack:ud-tmp", 0)

vfs.set_default_drive("root")
vfs.searchpath("root:/?.lua;root:/?/init.lua;/?.lua;/?/init.lua")

if os.getenv("CARBON_DEBUGMODE") ~= "true" then -- strip debug info
	kvstore._set("carbon:strip_internal_bytecode", true)
else
	print("Carbon running in debug mode, bytecode won't be optimized in the cache.")
end
