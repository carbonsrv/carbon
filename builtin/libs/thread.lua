-- Threads
-- Has a buttload of issues, mostly bindings. Gotta looooove broken shit.

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

function thread.rpcthread() -- not working, issues with binding or something .-.
	local chan = com.create()

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
		com.send(chan, msgpack.pack({
			["f"] = string.dump(fn),
			["args"] = args,
		}))
		return true
	end
	local function recieve()
		local res = com.receive(chan)
		return true, unpack(res)
	end

	thread.spawn(function()
		local function pushback(...)
			com.send(chan, msgpack.pack({...}))
		end
		while true do
			local args = {}
			local cmd = com.receive(chan)
			local f, err = loadstring(cmd.f)
			if err ~= nil then
				pushback(false, err)
			else
				pushback(pcall(f, unpack(cmd.args)))
			end
		end
	end, {
		["chan"] = chan,
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

function thread.kvstore() -- doesn't work either .-.
	local chan = com.create()

	thread.spawn(function()
		local store = {}
		while true do
			local suc, cmd = pcall(msgpack.unpack, com.receive(chan))
			if suc then
				if cmd.value then
					store[cmd.name] = cmd.value
					com.send(chan, msgpack.pack({value=true, error=nil}))
				else
					com.send(chan, msgpack.pack({value=store[cmd.name], error=nil}))
				end
			else
				com.send(chan, msgpack.pack({value=nil, error=cmd}))
			end
		end
	end, {
		["chan"] = chan
	})

	return function(name, value)
		if name then
			com.send(chan, msgpack.pack({["name"]=name, ["value"]=value}))
			local res = com.receive(chan)
			if res.error then
				return false, res.error
			else
				return true, res.value
			end
		end
	end
end

return thread
