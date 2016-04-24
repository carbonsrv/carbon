-- Main Glue

-- Add webroot to path and cpath. (OLD)
local webroot_path = var.root .. (var.root:match("/$") and "" or "/")

--package.path = webroot_path.."?.lua;"..webroot_path.."?/init.lua;"..package.path
package.cpath = webroot_path.."?.so;"..webroot_path.."loadall.so;"..package.cpath

-- Custom package loaders so that you can require the libraries built into Carbon.
local cache_do_cache_prefix = "carbon:do_cache:"
local cache_dont_cache_physfs = "carbon:dont_cache:physfs"
local cache_key_prefix = "carbon:lua_module:"
local cache_key_app = cache_key_prefix .. "app:" -- for custom stuff
local cache_key_app_location = cache_key_prefix .. "app_location:"
local cache_key_asset = cache_key_prefix .. "asset:"
local cache_key_asset_location = cache_key_prefix .. "asset_location:"
local cache_key_physfs = cache_key_prefix .. "physfs:"
local cache_key_physfs_location = cache_key_prefix .. "physfs_location:"

-- Load bc cache from kvstore
function loadcache(name)
	local modname = tostring(name):gsub("%.", "/")

	local f_bc = kvstore._get(cache_key_app..modname)
	if f_bc then
		local f, err = loadstring(f_bc, kvstore._get(cache_key_app_location..modname))
		if err then error(err, 0) end
		return f
	end

	local f_bc = kvstore._get(cache_key_asset..modname)
	if f_bc then
		local f, err = loadstring(f_bc, kvstore._get(cache_key_asset_location..modname))
		if err then error(err, 0) end
		return f
	end

	local f_bc = kvstore._get(cache_key_physfs..modname)
	if f_bc then
		local f, err = loadstring(f_bc, kvstore._get(cache_key_physfs_location..modname))
		if err then error(err, 0) end
		return f
	end
	return "\n\tno stored bytecode in kvstore under '"..cache_key_asset..modname.."'" ..
		"\n\tno stored bytecode in kvstore under '"..cache_key_physfs..modname.."'"
end

-- Load from compiled in /builtin/libs
local function loadasset_libs(name)
	local modname = tostring(name):gsub("%.", "/")
	local location = "libs/" .. modname .. ".lua"

	local src = carbon.glue(location)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location)
		if err then error(err, 0) end
		kvstore._set(cache_key_asset..modname, string.dump(f))
		kvstore._set(cache_key_asset_location..modname, location)
		return f
	end

	local location_init = "libs/" .. modname .. "/init.lua"
	local src = carbon.glue(location_init)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location)
		if err then error(err, 0) end
		kvstore._set(cache_key_asset..modname, string.dump(f))
		kvstore._set(cache_key_asset_location..modname, location_init)
		return f
	end
	return "\n\tno lib asset '/" .. location .. "' (not compiled in)" ..
		"\n\tno lib asset '/" .. location_init .. "' (not compiled in)"
end

-- Load from compiled in /builtin/3rdparty
local function loadasset_thirdparty(name)
	local modname = tostring(name):gsub("%.", "/")
	local location = "3rdparty/" .. modname .. ".lua"
	local src = carbon.glue(location)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location)
		if err then error(err, 0) end
		kvstore._set(cache_key_asset..modname, string.dump(f))
		kvstore._set(cache_key_asset_location..modname, location)
		return f
	end

	local location_init = "3rdparty/" .. modname .. "/init.lua"
	local src = carbon.glue(location_init)
	if src ~= "" then
		-- Compile and return the module
		local f, err = loadstring(src, location)
		if err then error(err, 0) end
		kvstore._set(cache_key_asset..modname, string.dump(f))
		kvstore._set(cache_key_asset_location..modname, location_init)
		return f
	end
	return "\n\tno thirdparty asset '/" .. location .. "' (not compiled in)"..
		"\n\tno thirdparty asset '/" .. location_init .. "' (not compiled in)"
end

-- Load from physfs and cache if not disabled for module
local function loadphysfs(name)
	local modname = tostring(name):gsub("%.", "/")
	local location = modname .. ".lua"
	local src, err1 = fs.readfile(location)
	if src then
		-- Compile and return the module
		local f, err = loadstring(src, location)
		if err then error(err, 0) end
		if kvstore._get(cache_do_cache_prefix..modname) ~= false and kvstore._get(cache_dont_cache_physfs) ~= true then
			kvstore._set(cache_key_physfs..modname, string.dump(f))
			kvstore._set(cache_key_physfs_location..modname, location)
		end
		return f
	end

	local location_init = modname .. "/init.lua"
	local src, err2 = fs.readfile(location_init)
	if src then
		-- Compile and return the module
		local f, err = loadstring(src, location)
		if err then error(err, 0) end
		return f
	end
	return "\n\tno file '" .. location .. "' in webroot"..
		"\n\tno file '" .. location_init .. "' in webroot"
end

-- Flush physfs cache
function carbon.flush_cache(name)
	local modname = tostring(name):gsub("%.", "/")
	kvstore._del(cache_key_physfs..modname)
	kvstore._del(cache_key_physfs_location..modname)
	package.loaded[modname] = nil
end
-- Set the module to not cache
function carbon.dont_cache(name)
	if name then
		local modname = tostring(name):gsub("%.", "/")
		kvstore._set(cache_do_cache_prefix..modname, false)
	else
		kvstore._set(cache_dont_cache_physfs, true)
	end
end

-- Install the loaders so that it's called just before the normal Lua loaders
table.insert(package.loaders, 2, loadcache)
table.insert(package.loaders, 3, loadasset_libs)
table.insert(package.loaders, 4, loadasset_thirdparty)
table.insert(package.loaders, 5, loadphysfs)

-- Load wrappers
-- LazyLoader! An automatic lazy loading generator.
function carbon.lazyload_mark(tablename, path)
	path = path or tablename
	local old = _G[tablename] or {}
	--print("Marking "..tablename.." to be lazily loaded.")
	local oldmt = getmetatable(old)
	setmetatable(old, {
		__index=function(t, key)
			--print("Lazy loaded "..tablename)
			setmetatable(t, oldmt or nil)
			local r = require(path)
			if r ~= true then
				t = r
			end
			return t[key]
		end
	})
	_G[tablename] = old
end

require("wrappers.globalwrappers")

local wrappers = {
	fs=false,
	"io",
	"os",
	"kvstore",
	"table",
	"encoding",
	"termbox",
	"debug",
	"exec",
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
