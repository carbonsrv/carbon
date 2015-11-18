-- Main glue

-- Load Libraries/Wrappers.
local function loadasset(name, location)
	_G[name] = assert(loadstring(carbon.glue(location)))()
end

loadasset("tags", "libs/tags.lua")
loadasset("thread", "libs/thread.lua")
loadasset("msgpack", "libs/MessagePack.lua")

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
