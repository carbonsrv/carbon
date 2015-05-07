-- Middleware helpers, not really useful with anything but the server init script
function mw.new(fn)
	code = ""
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
	r, err = mw.New(code)
	if err ~= nil then
		error(err)
	end
	return r
end
