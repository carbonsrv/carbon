-- physfs wrappers
physfs = {}
function physfs.mount(dir, mp, app)
	if dir then
		local err = carbon._physfs_mount(dir, mp or "/", app or false)
		if err then
			error(err, 0)
		end
	else
		error("No dir given.")
	end
end

function physfs.exists(file)
	if file then
		return carbon._physfs_exists(file)
	else
		error("No file given.")
	end
end

function physfs.isDir(path)
	if path then
		return carbon._physfs_isDir(path)
	else
		error("No path given.")
	end
end

function physfs.mkdir(path)
	if path then
		local err = carbon._physfs_mkdir(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function physfs.umount(path)
	if path then
		local err = carbon._physfs_umount(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function physfs.delete(path)
	if path then
		local err = carbon._physfs_delete(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function physfs.setWriteDir(path)
	if path then
		local err = carbon._physfs_setWriteDir(path)
		if err then
			error(err, 0)
		end
	else
		error("No path given.")
	end
end

function physfs.getWriteDir()
	return carbon._physfs_getWriteDir()
end

function physfs.list(path)
	if path then
		local list, err = carbon._physfs_list(path)
		if err then
			return nil, err
		end
		return luar.slice2table(list)
	else
		error("No path given.")
	end
end

function physfs.modtime(path)
	if path then
		local mt, err = carbon._physfs_modtime(path)
		if err then
			return nil, err
		end
		return mt
	else
		error("No path given.")
	end
end

function physfs.readfile(path)
	if path then
		local content, err = carbon._physfs_readfile(path)
		if err then
			return nil, err
		end
		return content
	else
		error("No path given.")
	end
end

function physfs.size(path)
	if path then
		local size, err = carbon._physfs_size(path)
		if err then
			return nil, err
		end
		return size
	else
		error("No path given.")
	end
end

return physfs
