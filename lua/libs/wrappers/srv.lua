-- srv wrapper
-- Wraps around carbon.srv and provides the usual http verbs: GET, POST, PUT, DELETE, PATCH, HEAD and OPTIONS. Also provides way to prefix all routing things.

local M = srv or {}

-- Track if used or not.
M.used = false

-- Prefix related stuff...
M._prefix = M._prefix or "" -- defaults to none.

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
M._vhosts = M._vhosts or {}
M.use_vhosts = M.use_vhosts or false
M.vhost = M.vhost or "***"

M.finish_handler = function(type, path, h)
	if M.use_vhosts then
		M._vhosts[type] = M._vhosts[type] or {}
		M._vhosts[type][path] = M._vhosts[type][path] or {}
		if type ~= "Use" then
			M._vhosts[type][path][M.vhost] = h
		else
			table.insert(M._vhosts[type][path][M.vhost], h)
		end
	else
		carbon.srv[type](path, h)
	end
end

function M.Begin_VHOSTS(cur)
	M.vhost = cur or M.vhost
	M.use_vhosts = true
end

function M.Finish_VHOSTS()
	if M.use_vhosts then
		M.use_vhosts = false
		for type, paths in pairs(M._vhosts) do
			for path, vhosts in pairs(paths) do
				if type == "Use" then
					for vhost, middlewares in pairs(vhosts) do
						for _, mw in pairs(middlewares) do
							carbon.srv.Use(mw.VHOST_Middleware({[vhost] = mw}))
						end
					end
				elseif type == "Default" then
					carbon.srv.NoRoute(mw.VHOST(vhosts))
				else
					carbon.srv[type](path, mw.VHOST(vhosts))
				end
			end
		end
		M._vhosts = {}
		M.vhost = "***"
	end
end

-- General things.
function M.Use(middleware)
	if tostring(middleware) == "gin.HandlerFunc" or tostring(middleware) == "func(*gin.Context)" then
		M.used = true
		M.finish_handler("Use", "", middleware)
		--carbon.srv.Use(middleware)
	else
		error("Middleware is not the valid type! (gin.HandlerFunc)")
	end
end

function M.DefaultRoute(handler)
	if tostring(handler) == "func(*gin.Context)" then
		M.used = true
		--carbon.srv.NoRoute(handler)
		M.finish_handler("Default", "", handler)
	else
		error("Invalid handler.")
	end
end

-- HTTP Verbs
function M.GET(pattern, ...)
	if type(pattern) == "string" then
		local handlers = {...}
		local h = {}
		for _, handler in pairs(handlers) do
			if tostring(handler) == "func(*gin.Context)" then
				M.used = true
				table.insert(h, handler)
			elseif type(handler) == "function" then
				M.used = true
				table.insert(h, mw.new(handler))
			else
				error("Invalid handler.")
			end
		end
		M.finish_handler("GET", M._prefix .. pattern, unpack(h))
	else
		error("Invalid pattern.")
	end
end

function M.POST(pattern, ...)
	if type(pattern) == "string" then
		local handlers = {...}
		local h = {}
		for _, handler in pairs(handlers) do
			if tostring(handler) == "func(*gin.Context)" then
				M.used = true
				table.insert(h, handler)
			elseif type(handler) == "function" then
				M.used = true
				table.insert(h, mw.new(handler))
			else
				error("Invalid handler.")
			end
		end
		M.finish_handler("POST", M._prefix .. pattern, unpack(h))
	else
		error("Invalid pattern.")
	end
end

function M.PUT(pattern, ...)
	if type(pattern) == "string" then
		local handlers = {...}
		local h = {}
		for _, handler in pairs(handlers) do
			if tostring(handler) == "func(*gin.Context)" then
				M.used = true
				table.insert(h, handler)
			elseif type(handler) == "function" then
				M.used = true
				table.insert(h, mw.new(handler))
			else
				error("Invalid handler.")
			end
		end
		M.finish_handler("PUT", M._prefix .. pattern, unpack(h))
	else
		error("Invalid pattern.")
	end
end

function M.DELETE(pattern, ...)
	if type(pattern) == "string" then
		local handlers = {...}
		local h = {}
		for _, handler in pairs(handlers) do
			if tostring(handler) == "func(*gin.Context)" then
				M.used = true
				table.insert(h, handler)
			elseif type(handler) == "function" then
				M.used = true
				table.insert(h, mw.new(handler))
			else
				error("Invalid handler.")
			end
		end
		M.finish_handler("DELETE", M._prefix .. pattern, unpack(h))
	else
		error("Invalid pattern.")
	end
end

function M.PATCH(pattern, ...)
	if type(pattern) == "string" then
		local handlers = {...}
		local h = {}
		for _, handler in pairs(handlers) do
			if tostring(handler) == "func(*gin.Context)" then
				M.used = true
				table.insert(h, handler)
			elseif type(handler) == "function" then
				M.used = true
				table.insert(h, mw.new(handler))
			else
				error("Invalid handler.")
			end
		end
		M.finish_handler("PATCH", M._prefix .. pattern, unpack(h))
	else
		error("Invalid pattern.")
	end
end

function M.HEAD(pattern, ...)
	if type(pattern) == "string" then
		local handlers = {...}
		local h = {}
		for _, handler in pairs(handlers) do
			if tostring(handler) == "func(*gin.Context)" then
				M.used = true
				table.insert(h, handler)
			elseif type(handler) == "function" then
				M.used = true
				table.insert(h, mw.new(handler))
			else
				error("Invalid handler.")
			end
		end
		M.finish_handler("HEAD", M._prefix .. pattern, unpack(h))
	else
		error("Invalid pattern.")
	end
end

function M.OPTIONS(pattern, ...)
	if type(pattern) == "string" then
		local handlers = {...}
		local h = {}
		for _, handler in pairs(handlers) do
			if tostring(handler) == "func(*gin.Context)" then
				M.used = true
				table.insert(h, handler)
			elseif type(handler) == "function" then
				M.used = true
				table.insert(h, mw.new(handler))
			else
				error("Invalid handler.")
			end
		end
		M.finish_handler("OPTIONS", M._prefix .. pattern, unpack(h))
	else
		error("Invalid pattern.")
	end
end

return M
