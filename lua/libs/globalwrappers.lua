-- Global wrappers

function io.list(path)
	local path = tostring(path) or "."
	local res, err = carbon._io_list(path)
	if err then
		return nil, err
	else
		return luar.slice2table(res)
	end
end

function syntaxhl(text, customcss)
	if customcss then
		return carbon._syntaxhl(text, tostring(customcss))
	else
		return carbon._syntaxhl(text, "")
	end
end
