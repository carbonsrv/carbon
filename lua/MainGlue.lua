-- Main glue

-- Custom package loader so that you can require the libraries
local function loadasset(name)
	local location = "libs/" .. tostring(name):gsub("/", ".") .. ".lua"
	local src = carbon.glue(location)
	if src then
		-- Compile and return the module
		return assert(loadstring(src, location))
	end
	return "\n\tno asset '/" .. location .. "' (not compiled in)"
end

-- Install the loader so that it's called just before the normal Lua loader
table.insert(package.loaders, 2, loadasset)

-- Load a few builtin libs.
require("tags")
thread = require("thread")

-- Specific Tags and Aliases
function link(url, opt)
	if opt then
		return tag"a"[{href=url, unpack(opt)}]
	else
		return tag"a"[{href=url}]
	end
end

function script(args)
	if type(args) == "table" then
		return tag"script"[args]
	elseif args ~= nil then
		return tag"script"(tostring(args))
	else
		return tag"script"
	end
end

function css(args)
	local a = {type="text/css"}
	if type(args) == "table" then
		for k,v in pairs(args) do
			a[k] = v
		end
	else
		return tag"style"[a](tostring(args))
	end
	return tag"style"[a]
end

function doctype()
	return tag"!DOCTYPE"[{"html"}]:force_open()
end

-- Wrappers
function syntaxhl(text, customcss)
	if customcss then
		return _syntaxhlfunc(text, tostring(customcss))
	else
		return _syntaxhlfunc(text, "")
	end
end
