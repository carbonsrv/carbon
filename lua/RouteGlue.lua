-- Middleware helpers, not really useful with anything but the server init script
function mw.new(fn, newstate)
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
		r, err = mw.DLR_NS(code)
	else
		r, err = mw.DLR_RUS(code)
	end
	if err ~= nil then
		error(err)
	end
	return r
end
