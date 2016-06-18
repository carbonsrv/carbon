-- Carbon VFS
-- modular simplistic virtual filesystem handling designed for carbon

local vfs = {}

-- method caching
local err, tos, strm, strf, strsub, strgsub, tconcat = error, tostring, string.match, string.find, string.sub, string.gsub, table.concat

-- Helpers:
-- Relative <-> Absolute Path conversion
vfs.helpers = {}

-- from http://lua-users.org/wiki/SplitJoin
local function strsplt(self, sSeparator, nMax, bRegexp)
	assert(sSeparator ~= '')
	assert(nMax == nil or nMax >= 1)

	local aRecord = {}

	if self:len() > 0 then
		local bPlain = not bRegexp
		nMax = nMax or -1

		local nField, nStart = 1, 1
		local nFirst,nLast = self:find(sSeparator, nStart, bPlain)
		while nFirst and nMax ~= 0 do
			aRecord[nField] = self:sub(nStart, nFirst-1)
			nField = nField+1
			nStart = nLast+1
			nFirst,nLast = self:find(sSeparator, nStart, bPlain)
			nMax = nMax-1
		end
		aRecord[nField] = self:sub(nStart)
	end

	return aRecord
end


local function e(msg)
	err("VFS: "..msg, 0)
end

local function is_relative(path)
	return strsub(path, 1, 1) ~= "/" -- that's probably not the best...
end

local function abspath(rel, base)
	if strsub(rel, 1, 1) == "/" then -- just in case rel is actually absolute
		return rel
	end
	if rel == "" or rel == "." then
		return base -- didn't change path
	end

	base = strgsub(base, "^/", "")
	local path = strsplt(base, "/")
	local pathn = #path
	local relpath = strsplt(rel, "/")
	for i=0, #relpath do
		local elm = relpath[i]
		if elm == "." then -- ignore
		elseif elm == ".." then -- step back
			if pathn > 0 then
				path[pathn] = nil
				pathn = pathn - 1
			end
		elseif elm ~= nil and elm ~= "" then
			pathn = pathn + 1
			path[pathn] = elm
		end
	end
	return "/"..tconcat(path, "/")
end

vfs.helpers.is_relative = is_relative
vfs.helpers.abspath = abspath

-- Backends:
-- Backends have an init function with arguments to support it.
-- These functions should set up the env for the other functions, like read.

-- You should at least provide read for a read-only fs, add write if you have a read-write fs. However, you should probably add exists, size, modtime, rename, mkdir, delete, isdir, getcwd, chdir and list if you want to have a fully fledged filesystem that is read-writable, has directories and keeps track of the current directory.

vfs.backends = {}

function vfs.backends.null() -- stubs, not actually usable
	local cwd = "/"
	return {
		write = function(loc, str) print("Write to "..abspath(loc, cwd).." with content: "..str) return true end,
		read = function(loc) print("Reading "..abspath(loc, cwd)) return "" end,
		size = function(loc) print("size check at "..abspath(loc, cwd)) return 0 end,
		exists = function(loc) print("check if "..abspath(loc, cwd).." exists") return true end,
		mkdir = function(loc) print("mkdir at "..abspath(loc)) return true end,
		delete = function(loc) print("delete at "..abspath(loc, cwd)) return true end,
		list = function(loc) print("list at "..abspath(loc, cwd)) return {} end,
		chdir = function(loc) print("cwd is now "..abspath(loc, cwd)) cwd = abspath(loc, cwd) return cwd end,
		getcwd = function() print("getcwd (is "..cwd..")") return cwd end,
	}
end

function vfs.backends.native(drivename, path) -- native backend
	local cwd = "/"
	local base = (path or (os.pwd and os.pwd()) or "") .. "/"

	local function getpath(rel)
		return base .. abspath(rel or ".", cwd)
	end


	local drv = { -- create all always existing funcs here
		read = function(loc)
			local fh, err = io.open(getpath(loc), "r")
			if err then
				return nil, err
			end
			local txt = fh:read("*all")
			fh:close()
			return txt
		end,
		write = function(loc, txt)
			local fh, err = io.open(getpath(loc), "w")
			if err then
				return nil, err
			end
			fh:write(txt)
			fh:close()
			return true
		end,
		delete = function(loc)
			return os.remove(getpath(loc))
		end,
		rename = function(loc1, loc2)
			return os.rename(getpath(loc1), getpath(loc2))
		end,
	}
	
	-- Mostly carbon specific additions.
	-- TODO: Maybe add LFS support?
	if os.exists then
		drv.exists = function(loc)
			return os.exists(getpath(loc))
		end
	end
	if io.list then
		drv.list = function(loc)
			return io.list(getpath(loc))
		end
	end
	if io.modtime then
		drv.modtime = function(loc)
			return io.modtime(getpath(loc))
		end
	end
	if io.isDir then
		drv.isdir = function(loc)
			return io.isDir(getpath(loc))
		end
	end
	if io.size then
		drv.size = function(loc)
			return io.size(getpath(loc))
		end
	end

	return drv
end

if carbon then
	local physfs = physfs or fs
	function vfs.backends.physfs(drivename, path, ismounted) -- read-only physfs backend for carbon
		if not ismounted then
			physfs.mount(path, "/"..drivename)
		end
		local cwd = "/"

		local base = "/"..drivename
		local function getdir(path)
			return base .. abspath(path or ".", cwd)
		end

		return {
			-- disabled modifying functions
			write = function() e("Drive "..drivename..": PhysFS writing disabled!") end,
			mkdir = function() e("Drive "..drivename..": PhysFS directory creation disabled!") end,
			delete = function() e("Drive "..drivename..": PhysFS file removal disabled!") end,
			rename = function() e("Drive "..drivename..": PhysFS renaming disabled!") end,

			-- read only funcs
			exists = function(loc) return physfs.exists(getdir(loc)) end,
			isdir = function(loc) return physfs.isDir(getdir(loc)) end,
			read = function(loc) return physfs.readfile(getdir(loc)) end,
			list = function(loc) return physfs.list(getdir(loc)) end,
			modtime = function(loc) return physfs.modtime(getdir(loc)) end,
			size = function(loc) return physfs.size(getdir(loc)) end,
	
			-- generic functions
			chdir = function(loc) cwd = abspath(loc, cwd) return cwd end,
			getcwd = function(loc) return cwd end,

			-- deinit function
			unmount = function() if not ismounted then physfs.unmount(base) end end,
		}
	end
end

-- drives:
-- Basically instances of backends

vfs.drives = {}

local function get_drive_field(drive, field)
	local drv = vfs.drives[drive]
	if drv then
		return drv[field]
	end
	e("No such drive: "..drive)
end

local function call_backend(drive, func, ...)
	local f = get_drive_field(drive, func)
	if f then
		return f(...)
	end
	e("Drive "..drive.." provides no function named "..func)
end

-- drive init and unmount
function vfs.new(drivename, backend, ...)
	if vfs.backends[backend] then
		if not vfs.drives[drivename] then
			vfs.drives[drivename] = vfs.backends[backend](drivename, ...)
			return true
		end
		return false
	end
	e("No such backend: "..backend)
end

function vfs.unmount(drivename)
	if vfs.drives[drivename] then
		vfs.drives[drivename].unmount()
		vfs.drives[drivename] = nil
	end
end

-- default drive selection:
-- if the path is not in the form of "drive:whatever", use the default drive.

-- set vfs.default_drive to "root" or whatever drive you want the default.

local function parse_path(path, default)
	local drive, filepath = string.match(path or "", "^(%w-):(.+)$")
	if drive and filepath then -- full vfs path
		return drive, filepath
	else -- "normal" path, like /bla
		return default or vfs.default_drive, path
	end
end

-- Generic function addition
-- Magic!

setmetatable(vfs, {__index=function(tbl, name)
	return function(filepath, ...)
		local drive, path = parse_path(filepath)
		if path then
			return call_backend(drive, name, path, ...)
		end
		return call_backend(drive, name, ...)
	end
end})


-- End
return vfs
