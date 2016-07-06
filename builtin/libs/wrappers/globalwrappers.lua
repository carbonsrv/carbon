-- Global wrappers

function syntaxhl(text, customcss)
	if customcss then
		return carbon._syntaxhl(text, tostring(customcss))
	else
		return carbon._syntaxhl(text, "")
	end
end

-- overload modulo for strings to allow formatting
getmetatable("").__mod = function(a, b)
	if not b then
		return a
	elseif type(b) == "table" then
		return string.format(a, unpack(b))
	else
		return string.format(a, b)
	end
end

-- Allow character access using overloading
getmetatable("").__index = function(s,k)
	if type(k)=="number" then
		return string.sub(s,k,k)
	else
		return string[k]
	end
end

-- Previous Carbon version had unixtime() as a special function since not on every single platform os.time() returns a unix timestamp.
-- However, most platforms do. Any POSIX-y one, at least.
unixtime = os.time
