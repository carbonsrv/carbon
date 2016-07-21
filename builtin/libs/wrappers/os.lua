-- os wrappers
function os.exists(path)
	if path then
		return carbon._os_exists(path)
	else
		error("No path given.")
	end
end

function os.chdir(path)
	if path then
		return carbon._os_chdir(path)
	end
	error("No path given.")
end

function os.abspath(path)
	if path then
		return assert(carbon._os_abspath(path))
	end
	error("No path given.")
end

function os.sleep(secs)
	if tonumber(secs) then
		carbon._os_sleep(tonumber(secs)*1000)
	end
	error("secs not a number!")
end

function os.pwd()
	local pwd, err = carbon._os_pwd()
	if err then
		error(err)
	end
	return pwd
end

function os.removeall(path)
	if path then
		return carbon._os_removeall(path)
	end
	error("No path given.")
end

function os.mkdir(path, perms)
	if path then
		return assert(carbon._os_mkdir(path, perms or 0755))
	end
	error("No path given.")
end
