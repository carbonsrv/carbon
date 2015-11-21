-- Main glue

-- Custom package loader so that you can require the libraries
local function loadasset(name)
	local location = "libs/" .. tostring(name):gsub("%.", "/") .. ".lua"
	local src = carbon.glue(location)
	if src ~= "" then
		-- Compile and return the module
		return assert(loadstring(src, location))
	end
	return "\n\tno asset '/" .. location .. "' (not compiled in)"
end

-- Install the loader so that it's called just before the normal Lua loaders
table.insert(package.loaders, 2, loadasset)

-- Load a few builtin libs.
require("globalwrappers")
require("tags")

thread = require("thread")
