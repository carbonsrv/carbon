-- io wrappers
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
