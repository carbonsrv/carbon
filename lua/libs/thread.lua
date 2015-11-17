-- Threads

function thread.spawn(fn, bindings)
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
		err = thread._spawn(code, true, bindings)
	else
		err = thread._spawn(code, false, {["s"]="v"})
	end
	if err ~= nil then
		return false, error(err)
	end
	return true
end
function thread.rpcthread() -- not working, issues with binding .-.
	local to_func = com.create()
	local to_len = com.create()
	local to_vals = com.createBuffered(64)
	local from_len = com.create()
	local from_vals = com.createBuffered(64)

	local function call(f, ...)
		local fn
		print(f)
		if type(f) == "function" then
			fn = fn
		elseif type(f) == "string" then
			local f, err = loadstring(f)
			if err then
				return false, err
			end
			fn = f
		end
		print(fn)
		local args = {...}
		print(com.send(to_func, string.dump(fn)))
		com.send(to_len, #args)
		for _, v in pairs(vars) do
			com.send(to_vals, v)
		end
		return true
	end
	local function recieve()
		local vals = {}
		local l = com.receive(from_len)
		if l > 0 then
			for i=1, l do
				vals[i] = com.receive(from_vals)
			end
		end
		return true, unpack(vals)
	end

	thread.spawn(function()
		print(luar.type(to_func))
		print(luar.type(to_len))
		local function pushback(...)
			local vars = {...}
			com.send(from_len, #vars)
			for _, v in pairs(vars) do
				com.send(from_vals, v)
			end
		end
		while true do
			local args = {}
			local bcode = com.receive(to_func)
			print(bcode)
			local l = com.receive(to_len)
			if l > 0 then
				for i=1, l do
					args[i] = com.receive(to_vals)
				end
			end
			local f, err = loadstring(bcode)
			if err ~= nil then
				pushback(false, err)
			else
				pushback(pcall(f, unpack(args)))
			end
		end
	end, {
		["to_func"] = to_func,
		["to_len"] = to_len,
		["to_vals"] = to_vals,
		["from_len"] = from_len,
		["from_vals"] = from_vals
	})
	return {
		["call_async"] = call,
		["call"] = (function(f, ...)
			local suc, err = call(f, ...)
			if not suc then
				return false, err
			end
			return recieve()
		end),
		["recieve"] = recieve,
	}
end

return threads
