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
	else
		error("No path given.")
	end
end

function os.abspath(path)
	if path then
		return assert(carbon._os_abspath(path))
	else
		error("No path given.")
	end
end

function os.sleep(secs)
	if tonumber(secs) then
		carbon._os_sleep(tonumber(secs)*1000)
	else
		error("secs not a number!")
	end
end

function os.pwd()
	local pwd, err = carbon._os_pwd()
	if err then
		error(err)
	end
	return pwd
end
