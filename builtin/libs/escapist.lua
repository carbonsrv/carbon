-- Various escape and unescape functions

local M = {}

local html_escape = {
	["<"]="&lt;",
	[">"]="&gt;",
	["&"]="&amp;"
}

local function uri_escape(a)
	return ("%%%02x"):format(a:byte())
end

local function uri_unescape(a)
	return string.char(tonumber(a,16))
end

M.escape = {
	html = function(str)
		return (str:gsub("[<>&]", html_escape))
	end,
	url = function(str)
		return (str:gsub("[^a-zA-Z0-9_.~-]", uri_escape))
	end,
	shell = function(str)
		return (str:gsub("[%s`~!#$&*()|\\'\";<>?{}[%]^]", "\\%1"))
	end,
	doublequotes = function(str)
		return (str:gsub("\\\"", "\\%1"))
	end
}

M.unescape = {
	url = function(str)
		return (str:gsub("+", " "):gsub("%%(%x%x)", uri_unescape))
	end
}

return M
