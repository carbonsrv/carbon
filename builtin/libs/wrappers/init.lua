-- Init Glue for Config State
kvstore.set("tmp:msgpack:ud-tmp", 0)

if os.getenv("CARBON_DEBUGMODE") ~= "true" then -- strip debug info
	kvstore._set("carbon:strip_internal_bytecode", true)
end
