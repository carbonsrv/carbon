--[[
Copyright &copy; 2015 Charles "ChickenNuggers" Heywood <vandor2012@gmail.com>
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
]]

local function serialize_value(value, indent, seen, co)
	local t = type(value)
	if t == "string" then
		return string.format("%q", value)
	elseif t == "number" or t == "boolean" or t == "nil" then
		return tostring(value)
	elseif t == "table" then
		if seen[value] then
			error("Recursive table, found " .. tostring(value) .. " repeated")
		else
			seen[value] = true
		end
		if not next(value) then
			-- Simple, return {}
			co("{}")
		else
			seen_now = {}
			-- Not so simple
			co("{")
			for i=1, #value do
				seen_now[i] = true
				co("\n" .. indent .. serialize_value(value[i], indent .. "\t", seen, co) .. ",")
			end
				-- For some reason, won't match numbers past a non-nil
				-- :: TODO ::
			for i, v in ipairs(value) do
				if not seen_now[i] then
					seen_now[i] = true
					co("\n" .. indent .. "[" .. tostring(i) .. "] = " .. serialize_value(value[i], indent .. "\t", seen, co) .. ",")
				end
			end
			-- Now, force [] around non-[A-Za-z_][0-9A-Za-z_]* names
			for k, v in pairs(value) do
				if not seen_now[k] then
					seen_now[k] = true
					if not tonumber(k) and k:match("^[A-Za-z_][0-9A-Za-z_]*$") then
						co("\n" .. indent .. k .. " = " .. serialize_value(value[k], indent .. "\t", seen, co) .. ",")
					elseif tonumber(k) then
						co("\n" .. indent .. "[" .. k .. "] = " .. serialize_value(value[k], indent .. "\t", seen, co) .. ",")
					else
						co("\n" .. indent .. "[" .. string.format("%q", k) .. "] = " .. serialize_value(value[k], indent .. "\t", seen, co) .. ",")
					end
				end
			end
			co("\n" .. indent:sub(2) .. "}")
		end
	else
		return tostring(value)
	end
	-- finished with no issues
	return co(true)
end

local function serialize(tbl, callback, is_coroutinable, ...)
	local clbk
	if not is_coroutinable then
		clbk = coroutine.wrap(function()
			while true do
				callback(coroutine.yield())
			end
		end)
	else
		clbk = coroutine.wrap(callback)
	end
	clbk(...) -- for cases such as passing files
	-- return in case coroutine yields false
	return serialize_value(tbl, "\t", {}, clbk)
end

local function stringify()
	-- coroutine to create string from a table
	local retvals = {}
	-- don't need to actually create a new string each time
	-- just call table.concat() at the end
	while true do
		local cur_val = coroutine.yield()
		if cur_val == true then
			break
		end
		table.insert(retvals, cur_val)
	end
	-- table.concat() returned values
	coroutine.yield(table.concat(retvals))
end

local function simple(tbl)
	return serialize(tbl, stringify, true)
end

local function tofile(file, close)
	while true do
		local cur_val = coroutine.yield()
		if cur_val == true then
			break
		end
		file:write(cur_val)
	end
	if close then
		file:close()
	end
end

return {
	stringify = stringify,
	serialize = serialize,
	simple = simple,
	file = tofile
}
