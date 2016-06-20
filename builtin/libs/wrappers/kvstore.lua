-- kvstore
function kvstore.set(name, value, timeout)
	if name then
		if value ~= nil then
			if timeout then
				kvstore._set_timeout(tostring(name), value, tonumber(timeout))
			end
			kvstore._set(tostring(name), value)
		else
			kvstore._del(tostring(name))
		end
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
		if not type(res) == "string" then
			local t = tostring(res)
			if t == "map[string]interface {}" then
				return luar.map2table(res)
			elseif t == "[]interface {}" then
				return luar.slice2table(res)
			end
		end
		return res
	else
		error("No name given.")
	end
end

function kvstore.inc(name, number)
	if name then
		kvstore._inc(tostring(name), tonumber(number) or 1)
	else
		error("No name given.")
	end
end

function kvstore.dec(name, number)
	if name then
		kvstore._dec(tostring(name), tonumber(number) or 1)
	else
		error("No name given.")
	end
end
