-- Return function
function content(data, code, ctype)
	local code = code or 200
	local ctype = ctype or "text/html"
	local content = ""
	if type(data) == "string" then
		content = data
	elseif type(data) == "table" and data.render ~= nil then
		content = data:render()
	else
		content = tostring(data)
	end
	context.Data(code, ctype, convert.stringtocharslice(content))
end

-- Vars and stuff from context.
function param(name)
	if name ~= nil then
		local f = _paramfunc(tostring(name))
		if f == "" then
			return nil
		end
		return f
	end
end
params = param
function form(name)
	if name ~= nil then
		local f = _formfunc(tostring(name))
		if f == "" then
			return nil
		end
		return f
	end
end
function query(name)
	if name ~= nil then
		local f = _queryfunc(tostring(name))
		if f == "" then
			return nil
		end
		return f
	end
end
