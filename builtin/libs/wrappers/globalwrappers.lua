-- Global wrappers

-- Not only wrappers, but also enhancements.
-- Pretty much the lua version of bind.go.

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

function io.glob(path)
	local path = path or "*"
	local res, err = carbon._io_glob(path)
	if err then
		return nil, err
	else
		return luar.slice2table(res)
	end
end

function io.modtime(path)
	local path = path or "*"
	local res, err = carbon._io_modtime(path)
	if err then
		return nil, err
	else
		return res
	end
end

function io.isDir(path)
	if path then
		return carbon._io_isDir(path)
	else
		error("No path given.")
	end
end

function io.size(path)
	if path then
		local size, err = carbon._io_size(path)
		if err then
			return nil, err
		end
		return size
	else
		error("No path given.")
	end
end

-- fs
fs = {}
function fs.mount(dir, mp, app)
	if dir then
		local err = carbon._fs_mount(dir, mp or "/", app or false)
		if err then
			error(err, 0)
		end
	else
		error("No dir given.")
	end
end

function fs.exists(file)
	if file then
		return carbon._fs_exists(file)
	else
		error("No file given.")
	end
end

function fs.isDir(path)
	if path then
		return carbon._fs_isDir(path)
	else
		error("No path given.")
	end
end

function fs.mkdir(path)
	if path then
		local err = carbon._fs_mkdir(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function fs.umount(path)
	if path then
		local err = carbon._fs_umount(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function fs.delete(path)
	if path then
		local err = carbon._fs_delete(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function fs.setWriteDir(path)
	if path then
		local err = carbon._fs_setWriteDir(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function fs.getWriteDir()
	return carbon._fs_getWriteDir()
end

function fs.list(path)
	if path then
		local list, err = carbon._fs_list(path)
		if err then
			return nil, err
		end
		return luar.slice2table(list)
	else
		error("No path given.")
	end
end

function fs.modtime(path)
	if path then
		local mt, err = carbon._fs_modtime(path)
		if err then
			return nil, err
		end
		return mt
	else
		error("No path given.")
	end
end

function fs.readfile(path)
	if path then
		local content, err = carbon._fs_readfile(path)
		if err then
			return nil, err
		end
		return content
	else
		error("No path given.")
	end
end

function fs.size(path)
	if path then
		local size, err = carbon._fs_size(path)
		if err then
			return nil, err
		end
		return size
	else
		error("No path given.")
	end
end

-- os
function os.exists(path)
	if path then
		return carbon._os_exists(path)
	else
		error("No path given.")
	end
end

function os.chdir(path)
	if path then
		return carbon._os_chdir(path)
	else
		error("No path given.")
	end
end

function os.abspath(path)
	if path then
		return assert(carbon._os_abspath(path))
	else
		error("No path given.")
	end
end

function os.sleep(secs)
	if tonumber(secs) then
		carbon._os_sleep(tonumber(secs))
	else
		error("secs not a number!")
	end
end

function os.pwd()
	local pwd, err = carbon._os_pwd()
	if err then
		error(err)
	end
	return pwd
end

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

-- table
table.unpack = unpack

function table.scopy(orig) -- http://lua-users.org/wiki/CopyTable
	local orig_type = type(orig)
	local copy
	if orig_type == 'table' then
		copy = {}
		for orig_key, orig_value in pairs(orig) do
			copy[orig_key] = orig_value
		end
	else -- number, string, boolean, etc
		copy = orig
	end
	return copy
end

function table.reverse(tbl)
	local len = #tbl
	local ret = {}

	for i = len, 1, -1 do
		ret[len - i + 1] = tbl[i]
	end

	return ret
end

function table.flip(tbl)
	local flipped = {}
	for k, v in pairs(tbl) do
		flipped[v] = k
	end
	return flipped
end

--[[
Ordered table iterator, allow to iterate on the natural order of the keys of a
table.

Example:
]]

local function __genOrderedIndex( t )
	local orderedIndex = {}
	for key in pairs(t) do
		table.insert( orderedIndex, key )
	end
	table.sort( orderedIndex )
	return orderedIndex
end

local function orderedNext(t, state)
	-- Equivalent of the next function, but returns the keys in the alphabetic
	-- order. We use a temporary ordered key table that is stored in the
	-- table being iterated.

	key = nil
	--print("orderedNext: state = "..tostring(state) )
	if state == nil then
		-- the first time, generate the index
		t.__orderedIndex = __genOrderedIndex( t )
		key = t.__orderedIndex[1]
	else
		-- fetch the next value
		for i = 1,table.getn(t.__orderedIndex) do
			if t.__orderedIndex[i] == state then
				key = t.__orderedIndex[i+1]
			end
		end
	end

	if key then
		return key, t[key]
	end

	-- no more value to return, cleanup
	t.__orderedIndex = nil
	return
end

function table.orderedPairs(t)
	-- Equivalent of the pairs() function on tables. Allows to iterate
	-- in order
	return orderedNext, t, nil
end

-- Similar to the above, but reverse!

local function __genReverseOrderedIndex( t )
	local orderedIndex = {}
	for key in pairs(t) do
		table.insert( orderedIndex, key )
	end
	table.sort( orderedIndex )
	return table.reverse(orderedIndex)
end

local function reverseOrderedNext(t, state)
	-- Equivalent of the next function, but returns the keys in the alphabetic
	-- order. We use a temporary ordered key table that is stored in the
	-- table being iterated.

	key = nil
	--print("orderedNext: state = "..tostring(state) )
	if state == nil then
		-- the first time, generate the index
		t.__orderedIndex = __genReverseOrderedIndex( t )
		key = t.__orderedIndex[1]
	else
		-- fetch the next value
		for i = 1,table.getn(t.__orderedIndex) do
			if t.__orderedIndex[i] == state then
				key = t.__orderedIndex[i+1]
			end
		end
	end

	if key then
		return key, t[key]
	end

	-- no more value to return, cleanup
	t.__orderedIndex = nil
	return
end

function table.reverseOrderedPairs(t)
	-- Equivalent of the pairs() function on tables. Allows to iterate
	-- in order
	return reverseOrderedNext, t, nil
end

-- encoding
encoding = {}
encoding.base64 = {}
function encoding.base64.decode(str)
	if str then
		local data, err = carbon._enc_base64_dec(str)
		if err ~= nil then
			error(err)
		end
		return data
	end
end
function encoding.base64.encode(str)
	if str then
		return carbon._enc_base64_enc(str)
	end
end

-- mime
mine = {}
function mime.byext(ext)
	return carbon._mime_byext(ext)
end

function mime.bytype(type)
	local exts, err = carbon._mime_bytype(type)
	if err then
		return nil, err
	end
	return luar.slice2table(exts)
end

-- debug
function debug.getallupvals(f)
	local i = 1
	local r = {}
	while true do
		local n, v = debug.getupvalue(f, i)
		if not n then
			if r == {} then
				return nil
			end
			return r
		end
		if n ~= "_ENV" then
			r[i] = {name=n, value=v}
		end
		i = i + 1
	end
end

function debug.setallupvals(f, vals)
	for i, pair in pairs(vals) do
		debug.setupvalue(f, i, pair["value"])
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
