-- kvstore
kvstore = kvstore or {}
function kvstore.set(name, value)
	if name then
		kvstore._set(tostring(name), value)
	else
		error("No name given.")
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

function kvstore.inc(name, number)
	if name then
		kvstore._inc(tostring(name), tonumber(number) or 1)
	else
		error("No name or number given.")
	end
end

function kvstore.dec(name, number)
	if name then
		kvstore._dec(tostring(name), tonumber(number) or 1)
	else
		error("No name given.")
	end
end
