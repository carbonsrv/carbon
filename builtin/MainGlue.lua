-- Main Glue

-- Add webroot to path and cpath. (OLD)
local webroot_path = var.root .. (var.root:match("/$") and "" or "/")

--package.path = webroot_path.."?.lua;"..webroot_path.."?/init.lua;"..package.path
package.cpath = webroot_path.."?.so;"..webroot_path.."loadall.so;"..package.cpath

-- Custom package loaders so that you can require the libraries built into Carbon.
local function loadasset_libs(name)
	local location = "libs/" .. tostring(name):gsub("%.", "/") .. ".lua"
	local src = carbon.glue(location)
	if src ~= "" then
		-- Compile and return the module
		return assert(loadstring(src, location))
	end

	local location_init = "libs/" .. tostring(name):gsub("%.", "/") .. "/init.lua"
	local src = carbon.glue(location_init)
	if src ~= "" then
		-- Compile and return the module
		return assert(loadstring(src, location_init))
	end
	return "\n\tno lib asset '/" .. location .. "' (not compiled in)\n\tno lib asset '/" .. location_init .. "' (not compiled in)"
end

local function loadasset_thirdparty(name)
	local location = "3rdparty/" .. tostring(name):gsub("%.", "/") .. ".lua"
	local src = carbon.glue(location)
	if src ~= "" then
		-- Compile and return the module
		return assert(loadstring(src, location))
	end

	local location_init = "3rdparty/" .. tostring(name):gsub("%.", "/") .. "/init.lua"
	local src = carbon.glue(location_init)
	if src ~= "" then
		-- Compile and return the module
		return assert(loadstring(src, location_init))
	end
	return "\n\tno thirdparty asset '/" .. location .. "' (not compiled in)\n\tno thirdparty asset '/" .. location_init .. "' (not compiled in)"
end

local function loadphysfs(name)
	local location = tostring(name):gsub("%.", "/") .. ".lua"
	local src, err1 = fs.readfile(location)
	if src then
		-- Compile and return the module
		return assert(loadstring(src, location))
	end

	local location_init = tostring(name):gsub("%.", "/") .. "/init.lua"
	local src, err2 = fs.readfile(location_init)
	if src then
		-- Compile and return the module
		return assert(loadstring(src, location_init))
	end
	return "\n\tno file '" .. location .. "' in webroot ("..err1..")\n\tno file '" .. location_init .. "' in webroot ("..err2..")"
end

-- Install the loaders so that it's called just before the normal Lua loaders
table.insert(package.loaders, 2, loadasset_libs)
table.insert(package.loaders, 3, loadasset_thirdparty)
table.insert(package.loaders, 4, loadphysfs)

-- Load a few builtin libs.
require("wrappers.globalwrappers")
require("thread")
require("tags")

thread = thread or require("thread")
