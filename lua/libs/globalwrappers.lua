-- Global wrappers

function syntaxhl(text, customcss)
	if customcss then
		return _syntaxhlfunc(text, tostring(customcss))
	else
		return _syntaxhlfunc(text, "")
	end
end
