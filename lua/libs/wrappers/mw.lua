-- Middleware Wrappers

function mw.new(fn, bindings, instances, newstate)
	local code = ""
	local instances = instances or -1
	if type(fn) == "function" then
		code = string.dump(fn)
	elseif type(fn) == "string" then
		fn, err = loadstring(code)
		if not err then
			code = string.dump(fn)
		else
			error(err)
		end
	end
	local r
	local err
	if newstate then
		if type(bindings) == "table" then
			r, err = mw.DLR_NS(code, true, bindings)
		else
			r, err = mw.DLR_NS(code, false, {["s"]="v"})
		end
	else
		if type(bindings) == "table" then
			r, err = mw.DLR_RUS(code, instances, true, bindings)
		else
			r, err = mw.DLR_RUS(code, instances, false, {["s"]="v"})
		end
	end
	if err ~= nil then
		error(err)
	end
	return r
end

function mw.ws(fn, bindings, instances)
	local code = ""
	if type(fn) == "function" then
		code = string.dump(fn)
	elseif type(fn) == "string" then
		fn, err = loadstring(code)
		if not err then
			code = string.dump(fn)
		else
			error(err)
		end
	end
	local r
	local err
	local instances = instances or -1
	if type(bindings) == "table" then
		r, err = mw.DLRWS_RUS(code, instances, true, bindings)
	else
		r, err = mw.DLRWS_RUS(code, instances, false, {["s"]="v"})
	end
	if err ~= nil then
		error(err)
	end
	return r
end

function mw.echo(code, resp)
	local resp = tonumber(resp) or 200
	if type(code) == "string" then
		return mw.Echo(resp, code)
	elseif type(code) == "table" then
		return mw.Echo(resp, code:render())
	else
		return mw.echoText(resp, tostring(code))
	end
end

function mw.echoText(text, resp)
	return mw.EchoText((tonumber(resp) or 200), text)
end

function mw.syntaxhl(text, resp)
	return mw.echo((tonumber(resp) or 200), syntaxhl(text))
end

function mw.static(path, prefix)
	local path = path or ""
	local prefix = prefix or ""
	return carbon._staticserve(path, prefix)
end

function mw.CGI(path, args, cwd)
	if type(path) == "string" then
		local args = args or {}
		args["DOCUMENT_ROOT"] = var.webroot
		local cwd = cwd or ""
		local preparedargs = {}
		for k, v in pairs(args) do
			table.insert(preparedargs, tostring(k).."="..tostring(v))
		end
		return carbon._mw_CGI(path, cwd, preparedargs)
	else
		error("path not string.")
	end
end
