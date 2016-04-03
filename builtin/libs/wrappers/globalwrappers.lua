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
