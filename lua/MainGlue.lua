-- Main glue

-- Load 3rdparty libraries
msgpack = assert(loadstring(carbon.glue("3rdparty/MessagePack.lua")))()

-- Tags
local html_escape={["<"]="&lt;",[">"]="&gt;",["&"]="&amp;"}

local uri_escape=function(a)
	return ("%%%02x"):format(a:byte())
end

local uri_unescape=function(a)
	return string.char(tonumber(a,16))
end

escape={
	html=function(str)
		return (str:gsub("[<>&]",html_escape))
	end,
	url=function(str)
		return (str:gsub("[^a-zA-Z0-9_.~-]",uri_escape))
	end,
	shell=function(str)
		return (str:gsub("[%s`~!#$&*()|\\'\";<>?{}[%]^]","\\%1"))
	end
}

unescape={
	url=function(str)
		return (str:gsub("+"," "):gsub("%%(%x%x)",uri_unescape))
	end
}

-- tag metatable
local tagmeth={
	render=function(self)
		local content
		if self.content then
			content={}
			for i,v in ipairs(self.content) do
				if type(v)=="string" then
					content[i]=v
				elseif type(v)=="number" then
					content[i]=tostring(v)
				else
					content[i]=v:render()
				end
			end
		end
		local options
		if self.options and next(self.options) then
			options={}
			for k,v in pairs(self.options) do
				if type(k)=="number" then
					table.insert(options,tostring(v))
				else
					table.insert(options,k.."=\""..tostring(v):gsub("\"","&quot;").."\"")
				end
			end
		end
		if self.fclose then
			if content then
				return table.concat(content).."</"..self.name..">"
			end
			return "</"..self.name..">"
		elseif self.fopen then
			local result
			if options then
				result="<"..self.name.." "..table.concat(options," ")..">"
			else
				result="<"..self.name..">"
			end
			if content then
				return result..table.concat(content)
			end
			return result
		else
			local result
			if options then
				result="<"..self.name.." "..table.concat(options," ")
			else
				result="<"..self.name
			end
			if content then
				return result..">"..table.concat(content).."</"..self.name..">"
			end
			return result.." />"
		end
	end,
	add_content=function(self,...)
		if not self.content then
			self.content={}
		end
		for i=1,select('#',...) do
			local value=select(i,...)
			if type(value)=="string" then
				if #value~=0 then
					table.insert(self.content,escape.html(value))
				end
			else
				table.insert(self.content,value)
			end
		end
		return self
	end,
	force_content=function(self,...)
		if not self.content then
			self.content={}
		end
		for i=1,select('#',...) do
			local value=select(i,...)
			if type(value)=="string" then
				if #value~=0 then
					table.insert(self.content,value)
				end
			else
				table.insert(self.content,value)
			end
		end
		return self
	end,
	add_options=function(self,tbl)
		if not self.options then
			self.options={}
		end
		for k,v in pairs(tbl) do
			if type(k)=="number" then
				table.insert(self.options,v)
			else
				self.options[k]=v
			end
		end
		return self
	end,
	clear_content=function(self)
		self.content=nil
		return self
	end,
	clear_options=function(self)
		self.options=nil
		return self
	end,
	set_content=function(self,...)
		self:clear_content()
		return self:add_content(...)
	end,
	set_options=function(self,tbl)
		self:clear_options()
		return self:add_options(tbl)
	end,
	force_open=function(self)
		self.fopen=true
		return self
	end,
	force_close=function(self)
		self.fclose=true
		return true
	end
}

local tagmt={
	__index=function(self,k)
		if type(k)=="table" then
			return self:add_options(k)
		else
			return tagmeth[k]
		end
	end,
	__call=tagmeth.add_content
}

function tag(name)
	return setmetatable({name=name},tagmt)
end

-- Specific Tags and Aliases
function link(url, opt)
	if opt then
		return tag"a"[{href=url, unpack(opt)}]
	else
		return tag"a"[{href=url}]
	end
end

function script(args)
	if type(args) == "table" then
		return tag"script"[args]
	elseif args ~= nil then
		return tag"script"(tostring(args))
	else
		return tag"script"
	end
end

function css(args)
	local a = {type="text/css"}
	if type(args) == "table" then
		for k,v in pairs(args) do
			a[k] = v
		end
	else
		return tag"style"[a](tostring(args))
	end
	return tag"style"[a]
end

function doctype()
	return tag"!DOCTYPE"[{"html"}]:force_open()
end

-- Wrappers
function syntaxhl(text, customcss)
	if customcss then
		return _syntaxhlfunc(text, tostring(customcss))
	else
		return _syntaxhlfunc(text, "")
	end
end

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
