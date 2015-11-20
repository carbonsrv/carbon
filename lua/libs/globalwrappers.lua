-- Global wrappers

-- io
function io.list(path)
	local path = path or "."
	local res, err = carbon._io_list(path)
	if err then
		return nil, err
	else
		return luar.slice2table(res)
	end
end

-- kvstore
function kvstore.set(name, value)
	if name and value then
		kvstore._set(tostring(name), value)
	else
		error("No name or value given.")
	end
end

function kvstore.del(name)
	if name then
		kvstore._del(tostring(name))
	else
		error("No name given.")
	end
end

function kvstore.get(name)
	if name then
		local res = kvstore._get(name)
		local t = tostring(res)
		if t == "map[string]interface {}" then
			return luar.map2table(res)
		elseif t == "[]interface {}" then
			return luar.slice2table(res)
		else
			return res
		end
	else
		error("No name given.")
	end
end

function kvstore._inc(name, number)
	if name and tonumber(number) then
		kvstore._inc(tostring(name), tonumber(value))
	else
		error("No name or number given.")
	end
end

function kvstore._dec(name, number)
	if name and tonumber(number) then
		kvstore._dec(tostring(name), tonumber(value))
	else
		error("No name or number given.")
	end
end

-- Global
function syntaxhl(text, customcss)
	if customcss then
		return carbon._syntaxhl(text, tostring(customcss))
	else
		return carbon._syntaxhl(text, "")
	end
end
