-- Main Glue

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


-- Install the loader so that it's called just before the normal Lua loaders
table.insert(package.loaders, 2, loadasset_libs)
table.insert(package.loaders, 3, loadasset_thirdparty)

-- Load a few builtin libs.
require("globalwrappers")
require("tags")

thread = thread or require("thread")
