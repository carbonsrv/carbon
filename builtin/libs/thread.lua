-- Threads
-- Has a buttload of issues, mostly bindings. Gotta looooove broken shit.

local msgpack = require("msgpack")

function thread.spawn(code, bindings, buffer, dontusemsgpack)
	local fn
	if type(code) == "function" then
		fn = code
	elseif type(code) == "string" then
		fn, err = loadstring(code)
		if err then
			error(err, 0)
		end
	end
	local buffer = tonumber(buffer) or -1

	if dontusemsgpack then
		local ch
		local err
		if type(bindings) == "table" then
			ch, err = thread._spawn(string.dump(fn), true, bindings, buffer)
		else
			ch, err = thread._spawn(string.dump(fn), false, {["s"]="v"}, buffer)
		end
		if err ~= nil then
			error(err)
		end
		return ch
	else
		mpfn = string.dump(function()
			--THREAD
			msgpack = require("msgpack")
			local args = msgpack.unpack(THREAD_ARGS)
			THREAD_ARGS=nil
			local binds = args["bindings"]
			if binds then
				for k, v in pairs(binds) do
					_G[k] = v
				end
			end
			local threadfn = args.fn
			args = nil
			threadfn()
		end)
		local mpbinds_raw, err = msgpack.pack({
			["fn"] = fn,
			["bindings"] = bindings
		})
		if err ~= nil then
			error(err)
		end
		mpbindings = {THREAD_ARGS=mpbinds_raw}
		local ch, err = thread._spawn(mpfn, true, mpbindings, buffer)
		if err ~= nil then
			error(err)
		end
		return ch
	end
end

-- Primitive interface to threads.
function thread.run(fn, ...)
	msgpack = require("msgpack")
	mpfn = string.dump(function()
		msgpack = require("msgpack")
		local args = msgpack.unpack(THREAD_ARGS)
		THREAD_ARGS=nil
		local thrcom = threadcom
		local threadfn = args.fn
		local thread_args = args.args
		args = nil
		com.send(thrcom, msgpack.pack({pcall(threadfn, unpack(thread_args, 1, thread_args.n))}))
	end)
	local mpbinds_raw, err = msgpack.pack({
		["fn"] = fn,
		["args"] = table.pack(...)
		})
	if err ~= nil then
		error(err)
	end
	mpbindings = {THREAD_ARGS=mpbinds_raw}
	local ch, err = thread._spawn(mpfn, true, mpbindings, -1)
	if err ~= nil then
		error(err)
	end
	local resgot = false
	return function()
		if not resgot then
			resgot = true
			local tmp = msgpack.unpack(com.receive(ch))
			local res = unpack(tmp, 1, tmp.n)
			tmp = nil
			ch = nil
			return res
		end
	end
end


function thread.rpcthread() -- not working, issues with binding or something .-.
	local rpc = thread.spawn(function()
		local msgpack = require("msgpack")
		while true do
			local src = com.receive(threadcom)
			local args = msgpack.unpack(src)
			local f, err = loadstring(args.f)
			if not err then
				com.send(threadcom, msgpack.pack({pcall(f, unpack(args.args, 1, args.args.n))}))
			else
				com.send(threadcom, msgpack.pack({false, err}))
			end
		end
	end)

	function call(f, ...)
		com.send(rpc, msgpack.pack({
			f = f,
			args = table.pack(...)
		}))
	end

	local function recieve()
		local tmp = msgpack.unpack(com.receive(rpc))
		return unpack(tmp, 1, tmp.n)
	end

	return {
		["call_async"] = call,
		["call"] = (function(f, ...)
			call(f, ...)
			return recieve()
		end),
		["recieve"] = recieve,
	}
end

function thread.kvstore() -- doesn't work either .-.
	local chan = thread.spawn(function()
		local msgpack = require("msgpack")
		local store = {}
		while true do
			local suc, cmd = pcall(msgpack.unpack, com.receive(threadcom))
			if suc then
				if cmd.value then
					store[cmd.name] = cmd.value
					com.send(threadcom, msgpack.pack({value=true, error=nil}))
				else
					com.send(threadcom, msgpack.pack({value=store[cmd.name], error=nil}))
				end
			else
				com.send(threadcom, msgpack.pack({value=nil, error=cmd}))
			end
		end
	end)

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
