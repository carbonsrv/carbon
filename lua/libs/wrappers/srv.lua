-- srv wrapper
-- Wraps around srv and provides the usual http verbs: GET, POST, PUT, DELETE, PATCH, HEAD and OPTIONS. Also provides way to prefix all routing things.

local M = {}
local args = {...}
local real_srv = args[1]

-- Prefix related stuff...
M._prefix = "" -- defaults to none.

function M.setPrefix(new_prefix)
	if type(new_prefix) == "string" then
		M._prefix = new_prefix
		return true
	else
		error("Prefix not a string!")
	end
end

function M.resetPrefix()
	M._prefix = ""
end

-- General things.
function M.Use(middleware)
	if tostring(middleware) == "gin.HandlerFunc" or tostring(middleware) == "func(*gin.Context)" then
		real_srv.Use(middleware)
	else
		error("Middleware is not the valid type! (gin.HandlerFunc)")
	end
end

function M.DefaultRoute(handler)
	if tostring(handler) == "func(*gin.Context)" then
		real_srv.NoRoute(handler)
	else
		error("Invalid handler.")
	end
end

-- Engine creation
function M.new()
	local s = carbon._gin_new()
	return require("wrappers.srv", s)
end

-- HTTP Verbs
function M.GET(pattern, handler, bindings)
	if type(pattern) == "string" then
		local h
		if tostring(handler) == "func(*gin.Context)" then
			h = handler
		elseif type(handler) == "function" then
			h = mw.new(handler, bindings)
		else
			error("Invalid handler.")
		end
		real_srv.GET(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

function M.POST(pattern, handler, bindings)
	if type(pattern) == "string" then
		local h
		if tostring(handler) == "func(*gin.Context)" then
			h = handler
		elseif type(handler) == "function" then
			h = mw.new(handler, bindings)
		else
			error("Invalid handler.")
		end
		real_srv.POST(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

function M.PUT(pattern, handler, bindings)
	if type(pattern) == "string" then
		local h
		if tostring(handler) == "func(*gin.Context)" then
			h = handler
		elseif type(handler) == "function" then
			h = mw.new(handler, bindings)
		else
			error("Invalid handler.")
		end
		real_srv.PUT(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

function M.DELETE(pattern, handler, bindings)
	if type(pattern) == "string" then
		local h
		if tostring(handler) == "func(*gin.Context)" then
			h = handler
		elseif type(handler) == "function" then
			h = mw.new(handler, bindings)
		else
			error("Invalid handler.")
		end
		real_srv.DELETE(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

function M.PATCH(pattern, handler, bindings)
	if type(pattern) == "string" then
		local h
		if tostring(handler) == "func(*gin.Context)" then
			h = handler
		elseif type(handler) == "function" then
			h = mw.new(handler, bindings)
		else
			error("Invalid handler.")
		end
		real_srv.PATCH(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

function M.HEAD(pattern, handler, bindings)
	if type(pattern) == "string" then
		local h
		if tostring(handler) == "func(*gin.Context)" then
			h = handler
		elseif type(handler) == "function" then
			h = mw.new(handler, bindings)
		else
			error("Invalid handler.")
		end
		real_srv.HEAD(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

function M.OPTIONS(pattern, handler, bindings)
	if type(pattern) == "string" then
		local h
		if tostring(handler) == "func(*gin.Context)" then
			h = handler
		elseif type(handler) == "function" then
			h = mw.new(handler, bindings)
		else
			error("Invalid handler.")
		end
		real_srv.OPTIONS(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

return M
