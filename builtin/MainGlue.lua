-- Main Glue

-- Add webroot to path and cpath. (OLD)
local webroot_path = var.root .. (var.root:match("/$") and "" or "/")

--package.path = webroot_path.."?.lua;"..webroot_path.."?/init.lua;"..package.path
package.cpath = webroot_path.."?.so;"..webroot_path.."loadall.so;"..package.cpath

-- Custom package loaders so that you can require the libraries built into Carbon.
local cache_do_cache_prefix = "carbon:do_cache:"
local cache_dont_cache_vfs = "carbon:dont_cache:vfs"
local cache_key_bc = "carbon:lua_module:bc:"
local cache_key_loc = "carbon:lua_module:loc:"

-- Load bc cache from kvstore
function loadcache(name)
	local modname = tostring(name):gsub("%.", "/")

	local f_bc = kvstore._get(cache_key_bc..modname)
	if f_bc then
		local f, err = loadstring(f_bc, kvstore._get(cache_key_loc..modname))
		if err then error(err, 0) end
		return f
	end

	return "\n\tno stored bytecode in kvstore under '"..cache_key_bc..modname.."'"
end

local function loadassets(name)
	-- Load from compiled in /builtin/libs
	local modname = tostring(name):gsub("%.", "/")
	local location_libs = "libs/" .. modname .. ".lua"
	local location_3rdparty = "3rdparty/" .. modname .. ".lua"
	local strip = kvstore._get("carbon:strip_internal_bytecode")

	local src = carbon.glue(location_libs)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location_libs)
		if err then error(err, 0) end
		kvstore._set(cache_key_bc..modname, string.dump(f, strip))
		kvstore._set(cache_key_loc..modname, location_libs)
		return f
	end

	local location_init_libs = "libs/" .. modname .. "/init.lua"
	local src = carbon.glue(location_init_libs)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location_libs)
		if err then error(err, 0) end
		kvstore._set(cache_key_bc..modname, string.dump(f, strip))
		kvstore._set(cache_key_loc..modname, location_init_libs)
		return f
	end

	-- Load from compiled in /builtin/3rdparty
	local src = carbon.glue(location_3rdparty)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location_3rdparty)
		if err then error(err, 0) end
		kvstore._set(cache_key_bc..modname, string.dump(f, strip))
		kvstore._set(cache_key_loc..modname, location_3rdparty)
		return f
	end

	local location_init_3rdparty = "3rdparty/" .. modname .. "/init.lua"
	local src = carbon.glue(location_init_3rdparty)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location_3rdparty)
		if err then error(err, 0) end
		kvstore._set(cache_key_bc..modname, string.dump(f, strip))
		kvstore._set(cache_key_loc..modname, location_init_3rdparty)
		return f
	end

	return "\n\tno lib asset '/" .. location_libs .. "' (not compiled in)" ..
		"\n\tno lib asset '/" .. location_init_libs .. "' (not compiled in)" ..
		"\n\tno thirdparty asset '/" .. location_3rdparty .. "' (not compiled in)"..
		"\n\tno thirdparty asset '/" .. location_init_3rdparty .. "' (not compiled in)"
end

-- Flush require cache
function carbon.flush_cache(name)
	local modname = tostring(name):gsub("%.", "/")
	kvstore._del(cache_key_bc..modname)
	kvstore._del(cache_key_loc..modname)
	package.loaded[modname] = nil
end
-- Set the module to not cache
function carbon.dont_cache(name)
	if name then
		local modname = tostring(name):gsub("%.", "/")
		kvstore._set(cache_do_cache_prefix..modname, false)
	else
		kvstore._set(cache_dont_cache_vfs, true)
	end
end

-- Install the loaders so that it's called just before the normal Lua loaders
table.insert(package.loaders, 2, loadcache)
table.insert(package.loaders, 3, loadassets)

-- Global wrappers
require("wrappers.globalwrappers")

-- VFS madness
require("wrappers.physfs")
vfs = require("vfs")
vfs.new("root", "physfs", nil, true)
table.insert(package.loaders, 4, vfs.loader)

-- Load wrappers
-- LazyLoader! An automatic lazy loading generator.
function carbon.lazyload_mark(tablename, path)
	path = path or tablename
	local old = _G[tablename] or {}
	--print("Marking "..tablename.." to be lazily loaded.")
	setmetatable(old, {
		__index=function(t, key)
			--print("Lazy loaded "..tablename)
			setmetatable(t, nil)
			local r = require(path)
			local res = {}
			if r ~= true then
				res = r
				_G[tablename] = r
			end
			return res[key] or (_G[tablename] or t)[key]
		end
	})
	_G[tablename] = old
end

local wrappers = {
	"fs",
	"io",
	"os",
	"kvstore",
	"table",
	"string",
	"encoding",
	"mime",
	"termbox",
	"debug",
	"exec",
	"sql",
}
for name, wrapper in pairs(wrappers) do
	if wrapper == false then
		require("wrappers."..name)
	elseif wrapper == true then
		carbon.lazyload_mark(name, "wrappers."..name)
	else
		carbon.lazyload_mark(wrapper, "wrappers."..wrapper)
	end
end

-- Load a few builtin libs.
carbon.lazyload_mark("thread")
require("tags")
