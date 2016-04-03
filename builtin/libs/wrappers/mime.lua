-- mime
mime = {}
function mime.byext(ext)
	local t = carbon._mime_byext(ext)
	if t == "" then
		return nil
	end
	return t
end

function mime.bytype(type)
	local exts, err = carbon._mime_bytype(type)
	if err then
		return nil, err
	end
	return luar.slice2table(exts)
end
