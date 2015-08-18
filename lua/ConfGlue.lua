-- Middleware helpers, not really useful with anything but the server init script
function mw.new(fn, bindings, newstate)
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
	if newstate then
		if type(bindings) == "table" then
			r, err = mw.DLR_NS(code, true, bindings)
		else
			r, err = mw.DLR_NS(code, false, {["s"]="v"})
		end
	else
		if type(bindings) == "table" then
			r, err = mw.DLR_RUS(code, true, bindings)
		else
			r, err = mw.DLR_RUS(code, false, {["s"]="v"})
		end
	end
	if err ~= nil then
		error(err)
	end
	return r
end
function mw.ws(fn, bindings)
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
	if type(bindings) == "table" then
		r, err = mw.DLRWS_RUS(code, true, bindings)
	else
		r, err = mw.DLRWS_RUS(code, false, {["s"]="v"})
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
