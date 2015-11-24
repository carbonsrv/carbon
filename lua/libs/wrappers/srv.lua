-- srv wrapper
-- Wraps around carbon.srv and provides the usual http verbs: GET, POST, PUT, DELETE, PATCH, HEAD and OPTIONS. Also provides way to prefix all routing things.

local M = {}

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

-- And VHOST stuff..
local finish_handler = function(type, path, h) carbon.srv[type](path, h) end
M._vhosts = {}
M.vhost = "***"

function M.Begin_VHOSTS(cur)
	M.vhost = cur or M.vhost
	finish_handler = function(type, path, h)
		M._vhosts[type] = M._vhosts[type] or {}
		M._vhosts[type][path] = M._vhosts[type][path] or {}
		M._vhosts[type][path][M.vhost] = h
	end
end

function M.Finish_VHOSTS()
	for type, paths in pairs(M._vhosts) do
		for path, vhosts in pairs(paths) do
			carbon.srv[type](path, mw.VHOST(vhosts))
		end
	end
end

-- General things.
function M.Use(middleware)
	if tostring(middleware) == "gin.HandlerFunc" or tostring(middleware) == "func(*gin.Context)" then
		carbon.srv.Use(middleware)
	else
		error("Middleware is not the valid type! (gin.HandlerFunc)")
	end
end

function M.DefaultRoute(handler)
	if tostring(handler) == "func(*gin.Context)" then
		carbon.srv.NoRoute(handler)
	else
		error("Invalid handler.")
	end
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
		finish_handler("GET", M._prefix .. pattern, h)
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
		finish_handler("POST", M._prefix .. pattern, h)
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
		finish_handler("PUT", M._prefix .. pattern, h)
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
		finish_handler("DELETE", M._prefix .. pattern, h)
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
		finish_handler("PATCH", M._prefix .. pattern, h)
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
		finish_handler("HEAD", M._prefix .. pattern, h)
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
		finish_handler("OPTIONS", M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

return M
