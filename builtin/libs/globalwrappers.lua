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

-- Global
function syntaxhl(text, customcss)
	if customcss then
		return carbon._syntaxhl(text, tostring(customcss))
	else
		return carbon._syntaxhl(text, "")
	end
end
