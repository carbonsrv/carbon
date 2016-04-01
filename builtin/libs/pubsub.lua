-- Com dispatcher/pubsub, publish to one, dispatched to many, similar to a mailing list!
-- It's a pubsub. :v

local _M = {}

thread = require("thread")
msgpack = require("msgpack")

_M.dispatcher = kvstore.get("pubsub:dispatcher_thread") or thread.spawn(function()
	local logger = require("libs.logger")
	local msgpack = require("msgpack")
	local thread = require("thread")

	local compaths = {}

	while true do
		local msg = msgpack.unpack(com.receive(threadcom))
		if msg.type == "sub" then
			local compath = compaths[msg.path]
			if not compath then
				compath = {}
			end
			table.insert(compath, msg.com)
			compaths[msg.path] = compath
			com.send(threadcom, nil)
		elseif msg.type == "unsub" then
			local compath = compaths[msg.path]
			if compath then
				for k, v in pairs(compath) do
					if v == msg.com then
						compaths[msg.path][k] = nil
						break
					end
				end
			end
		elseif msg.type == "pub" then
			local path = msg.path
			local message = msg.msg
			local sender = function()
				local compath = compaths[msg.path]
				if compath then
					for i, chan in pairs(compath) do
						if msg.try then
							com.try_send(chan, message) -- on mediocre msg spam/fast processors
						else
							com.send(chan, message)
						end
					end
				end
			end
			if msg.threaded then
				thread.spawn(sender) -- not quite as effective with the memory, but it will make the thinger not get stuck because of a dead chan.
			else
				sender()
			end
		end
	end
end)
kvstore.set("pubsub:dispatcher_thread", _M.dispatcher)

-- Subscribe to topic.
function _M.sub(path, chan, bindings, buffer)
	if chan then
		local chan = chan
		if type(chan) == "function" then
			chan = thread.spawn(chan, bindings, buffer or 64)
		end
		com.send(_M.dispatcher, msgpack.pack{
			type="sub",
			com=chan,
			path=path
		})
		com.receive(_M.dispatcher) -- Block until it's done for safety reasons.
	else
		error("chan not given!")
	end
end
function _M.unsub(path, chan)
	if chan then
		com.send(_M.dispatcher, msgpack.pack{
			type="unsub",
			com=chan,
			path=path
		})
		com.send(chan, nil)
		com.receive(_M.dispatcher) -- Block until it's done for safety reasons.
	else
		error("chan not given!")
	end
end

function _M.pub(path, msg, opt)
	local mustdeliver, send_threaded
	if opt == "threaded" then
		send_threaded = true
		mustdeliver = true
	elseif opt then
		send_threaded = false
		mustdeliver = true
	else
		send_threaded = false
		mustdeliver = false
	end
	if msg == nil then
		return false, "Can't publish nil value."
	end
	com.send(_M.dispatcher, msgpack.pack{
		type="pub",
		path=path,
		msg=msg,
		try=not mustdeliver,
		threaded=send_threaded,
	})
	return true
end

-- LTN12 compatibility helpers
function _M.subscriber(name) -- ltn12 compatible subscriber source
	local retcom = com.create()
	_M.sub(name, retcom)
	return function()
		if retcom then
			local src = com.receive(retcom)
			if src then
				return src
			else
				retcom = nil -- just to help the GC in case this function is kept somewhere
				return nil
			end
		end
	end
end

function _M.publisher(name) -- ltn12 compatible publisher sink
	return function(chunk)
		return _M.pub(name, chunk)
	end
end

return _M
