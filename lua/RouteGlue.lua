-- Middleware helpers, not really useful with anything but the server init script
function mw.new(fn)
	code = ""
	if type(fn) == "function" then
		code = string.dump(fn)
	elseif type(fn) == "string" then
		code = fn
	end
	r, err = mw.New(code)
	if err ~= nil then
		error(err)
	end
	return r
end
