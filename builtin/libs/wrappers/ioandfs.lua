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
