-- fs compatability layer for vfs
-- Carbon 1.3 added vfs, a more generic virtual filesystem framework. fs, which is now just a wrapper around vfs is still there for convenience, however some less used functions are stubs/unsupported.

fs = {}

-- Manual incompatability notices
function fs.mount()
	error("Please use vfs' vfs.new('drivename', 'physfs', '/path/to/archive.zip') or something similar instead of fs.mount.", 0)
end

function fs.unmount()
	error("Please use vfs' vfs.umount('drivename') instead of fs.unmount.", 0)
end

function fs.setWriteDir()
	print("Warning: Writing to physfs is no longer supported, since everything moved to vfs and no support exists for the physfs backend. Therefore, fs.setWriteDir is useless.")
end

function fs.getWriteDir()
	error("fs.getWriteDir got removed, since fs.setWriteDir is defunct.", 0)
end

-- Automated direct conversions.
local direct_conversions = {
	exists = "exists",
	isDir = "isdir",
	mkdir = "mkdir",
	delete = "delete",
	modtime = "modtime",
	readfile = "read",
	list = "list",
	size = "size",
}

for old, new in pairs(direct_conversions) do
	fs[old] = function(file, ...)
		return vfs[new]("root:"..(file or "/"), ...)
	end
end
