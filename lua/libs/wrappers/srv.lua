-- srv wrapper
-- Wraps around carbon.srv and provides the usual http verbs: GET, POST, PUT, DELETE, PATCH, HEAD and OPTIONS. Also provides way to prefix all routing things.

local M = srv or {}

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
	M.use_vhosts = false
	for type, paths in pairs(M._vhosts) do
		for path, vhosts in pairs(paths) do
			if type == "Use" then
				carbon.srv.Use(mw.VHOST_Middleware(vhosts))
			elseif type == "Default" then
				if #vhosts == 1 then
					local _, h = pairs(vhosts)(vhosts)
					carbon.srv.NoRoute(h)
				else
					carbon.srv.NoRoute(mw.VHOST(vhosts))
				end
			else
				if #vhosts == 1 then
					local _, h = pairs(vhosts)(vhosts)
					carbon.srv[type](path, h)
				else
					carbon.srv[type](path, mw.VHOST(vhosts))
				end
			end
		end
	end
	M._vhosts = {}
	M.vhost = "***"
end

-- General things.
function M.Use(middleware)
	if tostring(middleware) == "gin.HandlerFunc" or tostring(middleware) == "func(*gin.Context)" then
		--M.finish_handler("Use", M._prefix, middleware)
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
		M.finish_handler("GET", M._prefix .. pattern, h)
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
		M.finish_handler("POST", M._prefix .. pattern, h)
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
		M.finish_handler("PUT", M._prefix .. pattern, h)
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
		M.finish_handler("DELETE", M._prefix .. pattern, h)
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
		M.finish_handler("PATCH", M._prefix .. pattern, h)
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
		M.finish_handler("HEAD", M._prefix .. pattern, h)
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
		M.finish_handler("OPTIONS", M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

return M
