-- srv wrapper
-- Wraps around carbon.srv and provides the usual http verbs: GET, POST, PUT, DELETE, PATCH, HEAD and OPTIONS. Also provides way to prefix all routing things.

local M = {}

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
		carbon.srv.GET(M._prefix .. pattern, h)
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
		carbon.srv.POST(M._prefix .. pattern, h)
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
		carbon.srv.PUT(M._prefix .. pattern, h)
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
		carbon.srv.DELETE(M._prefix .. pattern, h)
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
		carbon.srv.PATCH(M._prefix .. pattern, h)
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
		carbon.srv.HEAD(M._prefix .. pattern, h)
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
		carbon.srv.OPTIONS(M._prefix .. pattern, h)
	else
		error("Invalid pattern.")
	end
end

return M
