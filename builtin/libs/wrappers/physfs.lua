-- physfs wrappers
physfs = {}
function physfs.mount(dir, mp, app)
	if dir then
		local err = carbon._physfs_mount(dir, mp or "/", app or false)
		if err then
			error(err, 0)
		end
		return true
	end
	error("physfs: mount: No dir given.")
end

function physfs.exists(file)
	if file then
		return carbon._physfs_exists(file)
	end
	error("physfs: exists: No file given.")
end

function physfs.isDir(path)
	if path then
		return carbon._physfs_isDir(path)
	end
	error("physfs: isDir: No path given.")
end

function physfs.mkdir(path)
	if path then
		local err = carbon._physfs_mkdir(path)
		if err then
			error(err, 0)
		end
		return true
	end
	error("physfs: mkdir: No path given.")
end

function physfs.umount(path)
	if path then
		local err = carbon._physfs_umount(path)
		if err then
			return err
		end
		return true
	end
	error("physfs: umount: No path given.")
end

function physfs.delete(path)
	if path then
		local err = carbon._physfs_delete(path)
		if err then
			error(err, 0)
		end
		return true
	end
	error("physfs: delete: No path given.")
end

function physfs.setWriteDir(path)
	if path then
		local err = carbon._physfs_setWriteDir(path)
		if err then
			error(err, 0)
		end
		return true
	end
	error("physfs: setWriteDir: No path given.")
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
	end
	error("physfs: list: No path given.")
end

function physfs.modtime(path)
	if path then
		local mt, err = carbon._physfs_modtime(path)
		if err then
			return nil, err
		end
		return mt
	end
	error("physfs: modtime: No path given.")
end

function physfs.readfile(path)
	if path then
		local content, err = carbon._physfs_readfile(path)
		if err then
			return nil, err
		end
		return content
	end
	error("physfs: readfile: No path given.")
end

function physfs.readat(path, at, count)
	if path and at and count then
		local s, err, count = carbon._physfs_readat(path, at, count)
		if err then
			return nil, err
		end
		return s, nil, count
	end
	error("physfs: readat: No path given.")
end
function physfs.readn(path, count)
	if path and count then
		local s, err, count = carbon._physfs_readn(path, count)
		if err then
			return nil, err
		end
		return s, nil, count
	end
	error("physfs: readn: No path given.")
end

function physfs.size(path)
	if path then
		local size, err = carbon._physfs_size(path)
		if err then
			return nil, err
		end
		return size
	end
	error("physfs: size: No path given.")
end

function physfs.needfile(path)
	if path then
		return carbon._physfs_needfile(path)
	end
	error("physfs: needfile: No path given.")
end

function physfs.close(path)
	if path then
		return carbon._physfs_close(path)
	end
	error("physfs: close: No path given.")
end

return physfs
