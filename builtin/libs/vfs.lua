-- Carbon VFS
-- modular simplistic virtual filesystem handling designed for carbon

vfs = {}

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
	error("VFS: "..msg, 0)
end

local function is_relative(path)
	return string.sub(path, 1, 1) ~= "/" -- that's probably not the best...
end

local function abspath(rel, base)
	if string.sub(rel, 1, 1) == "/" then -- just in case rel is actually absolute
		return rel
	end
	if rel == "" or rel == "." then
		return base -- didn't change path
	end

	base = string.gsub(base, "^/", "")
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
	return "/"..table.concat(path, "/")
end
vfs.abspath = abspath
vfs.is_relative = is_relative

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
		if not ismounted and not physfs.exists("/"..drivename) then
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

	-- Some very cool special goodie: A thread-safe proxy backend.
	-- Allows a single instance of a backend to be used by serveral threads, sharing it's cwd and whatnot.
	function vfs.backends.shared(drivename, sharedbackend, ...)
		local msgpack = require("msgpack")
		local kvstore_key_base = "carbon:vfs:"..drivename..":"
		if not sharedbackend then
			if not kvstore._get(kvstore_key_base.."com") then
				e("Shared backend has not been initialized for drive "..drivename)
			end

			local function call(name, ...)
				local shrd = kvstore._get(kvstore_key_base.."com")
				com.send(shrd, msgpack.pack({
					method = name,
					args = table.pack(...)
				}))

				local res = msgpack.unpack(com.receive(shrd))
				if res[1] == false then
					error(res[2], 0)
				elseif res[1] == true then
					return unpack(res, 2, res.n)
				end
			end
			return setmetatable({
				unmount = function(...)
					call("unmount")
					kvstore._del(kvstore_key_base.."com")
					kvstore._del(kvstore_key_base.."args")
				end
			}, {__index = function(_, name)
				return function(...) return call(name, ...) end
			end})
		else -- init backend and put the com in the kvstore
			local vfs_backend = vfs.backends[sharedbackend]
			kvstore._set(kvstore_key_base.."args", msgpack.pack(table.pack(...)))
			local shrd = thread.spawn(function()
				local msgpack = require("msgpack")
				local bargs = msgpack.unpack(kvstore._get(kvstore_key_base.."args"))
				vfs = require("vfs")
				local backend = vfs_backend(drivename, unpack(bargs, 1, bargs.n))

				while true do
					local src = com.receive(threadcom)
					local cmd = msgpack.unpack(src)

					local name, args = cmd.method, cmd.args
					local res
					if not backend[name] then
						res = {false, "VFS: Backend "..sharedbackend.." provides no function "..name}
					else
						res = table.pack(pcall(backend[name], unpack(args, 1, args.n)))
					end
					com.send(threadcom, msgpack.pack(res))
					if name == "unmount" then -- I guess it is time to go.
						return
					end
				end
			end)
			kvstore._set(kvstore_key_base.."com", shrd)
			return vfs.backends.shared(drivename)
		end
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

vfs.get_drive_field = get_drive_field
vfs.call_backend = call_backend

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

-- call vfs.default_drive with"root" or whatever drive you want the default.
local default_drive_key = "carbon:vfs:default_drive"
function vfs.set_default_drive(drivename)
	kvstore._set(default_drive_key, drivename)
end
function vfs.default_drive()
	return kvstore._get(default_drive_key)
end

local function parse_path(path, default)
	local drive, filepath = string.match(path or "", "^(%w-):(.+)$")
	if drive and filepath then -- full vfs path
		return drive, filepath
	else -- "normal" path, like /bla
		return default or vfs.default_drive(), path
	end
end
vfs.parse_path = parse_path

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
