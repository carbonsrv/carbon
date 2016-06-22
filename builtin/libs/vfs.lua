-- Carbon VFS
-- modular simplistic virtual filesystem handling designed for carbon

vfs = {}

-- Depends on: ltn12
local ltn12 = require("ltn12")

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
-- You can also implement other methods for your backends specific features.
-- One thing to implement would be LTN12 compatible reader and writer, however they don't have to exist.
-- If they don't, rather crude LTN12 wrappers around read and write exists as a fallback if they don't.
-- LTN12 compatibility would make things more efficient given proper implementation, because of the streaming ability, resulting in possibly faster transfer but mostly less memory use.

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
		reader = function(loc)
			return ltn12.source.file(io.open(getpath(loc)))
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
		writer = function(loc)
			return ltn12.sink.file(io.open(getpath(loc)))
		end,
		copy = function(src, dst)
			return ltn12.pump.all(
				ltn12.source.file(io.open(getpath(src))),
				ltn12.sink.file(io.open(getpath(dst)))
			)
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
	-- Read-only (for now) physfs backend for carbon
	-- TODO: Maybe write compatibility?
	local physfs = physfs or fs
	function vfs.backends.physfs(drivename, path, ismounted)
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
			reader = function(loc)
				local i = 1
				return function()
					local chunk, err = physfs.readat(getdir(loc), i, ltn12.BLOCKSIZE)
					i = i + ltn12.BLOCKSIZE
					if err or chunk == "" then
						return nil
					end
					return chunk, i
				end
			end,
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

	function vfs.backends.gofs(drivename, fs, prefix)
		if not fs then
			e("Backend GOFS: Needs http.FileSystem.")
		end

		local cwd = "/"
		local base = "/"..(prefix or "")
		local function getdir(path)
			return base .. abspath(path or ".", cwd)
		end

		return {
			-- disabled modifying functions
			write = function() e("Drive "..drivename..": writing disabled!") end,
			mkdir = function() e("Drive "..drivename..": directory creation disabled!") end,
			delete = function() e("Drive "..drivename..": file removal disabled!") end,
			rename = function() e("Drive "..drivename..": renaming disabled!") end,

			-- read only funcs
			exists = function(loc) return carbon._filesystem_exists(fs, getdir(loc)) end,
			isdir = function(loc) return carbon._filesystem_isdir(fs, getdir(loc)) end,
			read = function(loc)
				local str, err = carbon._filesystem_readfile(fs, getdir(loc))
				if err then
					return nil, err
				end
				return str
			end,
			reader = function(loc)
				local i = 1
				return function()
					local chunk, err = carbon._filesystem_readat(fs, getdir(loc), i, ltn12.BLOCKSIZE)
					i = i + ltn12.BLOCKSIZE
					if err or chunk == "" then
						return nil
					end
					return chunk, i
				end
			end,
			list = function(loc)
				local res, err = carbon._filesystem_list(fs, getdir(loc))
				if err then
					return nil, err
				end
				return luar.slice2table(res), nil
			end,
			modtime = function(loc)
				local res, err = carbon._filesystem_modtime(fs, getdir(loc))
				if err then
					return nil, err
				end
				return res, nil
			end,
			size = function(loc)
				local res, err = carbon._filesystem_size(fs, getdir(loc))
				if err then
					return nil, err
				end
				return res, nil
			end,

			-- generic functions
			chdir = function(loc) cwd = abspath(loc, cwd) return cwd end,
			getcwd = function(loc) return cwd end,
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
			local tmp = {
				unmount = function(...)
					call("unmount")
					kvstore._del(kvstore_key_base.."com")
					kvstore._del(kvstore_key_base.."args")
				end
			}
			for k, v in pairs(luar.slice2table(kvstore._get(kvstore_key_base.."methods"))) do
				if not k == "unmount" then 
					tmp[v] = function(...)
						return call(v, ...)
					end
				end
			end
			return tmp
		else -- init backend and put the com in the kvstore
			local vfs_backend = vfs.backends[sharedbackend]
			kvstore._set(kvstore_key_base.."args", msgpack.pack(table.pack(...)))
			local shrd = thread.spawn(function()
				local msgpack = require("msgpack")
				local bargs = msgpack.unpack(kvstore._get(kvstore_key_base.."args"))
				vfs = require("vfs")
				local drive = vfs_backend(drivename, unpack(bargs, 1, bargs.n))

				local methodlist = {}
				for k, _ in pairs(drive) do
					table.insert(methodlist, k)
				end
				kvstore._set(kvstore_key_base.."methods", methodlist)

				com.send(threadcom, true) -- indicate that we are done initializing

				while true do
					local src = com.receive(threadcom)
					local cmd = msgpack.unpack(src)

					local name, args = cmd.method, cmd.args
					local res = table.pack(pcall(drive[name], unpack(args, 1, args.n)))
					com.send(threadcom, msgpack.pack(res))
					if name == "unmount" then -- I guess it is time to go.
						return
					end
				end
			end)
			com.receive(shrd) -- block until thread is done with init
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

-- Helpers which overwrite backend functions.
function vfs.copy(fpsrc, fpdst)
	local fpsrc_drivename, fpsrc_path = parse_path(fpsrc)
	local fpdst_drivename, fpdst_path = parse_path(fpdst)
	local fpsrc_drive = vfs.drives[fpsrc]
	local fpdst_drive = vfs.drives[fpdst]
	local fpsrc_backend = vfs.backends[fpsrc_drive]
	local fpdst_backend = vfs.backends[fpdst_drive]

	if fpsrc_drive == fpdst_drive then -- same device copy
		if fpsrc_drive.copy then -- backend has a specific copy function
			return fpsrc_drive.copy(fpsrc_path, fpdst_path)
		end
	end

	-- streaming inter device copy
	if fpsrc_drive.reader and fpdst_drive.writer then -- ltn12! woo!
		return ltn12.pump.all(
			fpsrc_backend.reader(fpsrc_path),
			fpdst_backend.writer(fpdst_path)
		)
	end

	-- fallback
	local src = fpsrc_drive.read(fpsrc_path)
	return fpdst.write(fpdst_path, src)
end

function vfs.reader(src)
	local src_drivename, src_path = parse_path(src)
	local src_drive = vfs.drives[src_drivename]

	if src_drive.reader then -- native LTN12 reader exists
		return src_drive.reader(src_path)
	end

	-- A quite dirty hack, just for programs which definitly want a LTN12 reader, but the backend doesn't have one.
	-- Probably not that efficient.
	local src = src_drive.read(src_path)
	return ltn12.source.string(src)
end

function vfs.writer(dst)
	local dst_drivename, dst_path = parse_path(dst)
	local dst_drive = vfs.drives[dst_drivename]

	if dst_drive.writer then
		return dst_drive.writer(dst_path)
	end

	-- Another dirty hack. This one just concatenates chunk after chunk and writes it out when there is no more.
	local s = ""
	return function(chunk)
		if not chunk then
			return dst_drive.write(dst_path, chunk)
		end
		s = s .. tostring(chunk)
		return 1
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
