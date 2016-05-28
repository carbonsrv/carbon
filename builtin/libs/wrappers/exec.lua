-- exec wrappers

--Default read size
exec.READ_SIZE = 1

-- Convenience wrappers
function exec.source(...)
	local cmd = exec.exec(...)
	return exec.source_using(cmd)
end

function exec.filter(...)
	local cmd = exec.exec(...)
	return exec.filter_using(cmd)
end

function exec.sink(...)
	local cmd = exec.exec(...)
	return exec.sink_using(cmd)
end

-- Actual stuff.
function exec.source_using(cmd) -- LTN12 compatible exec call! \o/
	return function()
		local line, err = cmd.Read_Stdout(exec.READ_SIZE)
		if not err then
			return line
		else
			cmd.Close()
		end
	end
end

function exec.filter_using(cmd) -- returns two ltn12 functions, a sink and a source respectively.
	local sink = function(chunk)
		if chunk then
			local err = cmd.Write_Stdin(tostring(chunk))
			if not err then
				return 1
			else
				cmd.Close()
			end
		end
	end
	local source = function()
		local line, err = cmd.Read_Stdout(exec.READ_SIZE)
		if not err then
			return line
		end
	end
	return sink, source
end

function exec.sink_using(cmd)
	return function(chunk)
		if chunk then
			local err = cmd.Write_Stdin(tostring(chunk))
			if not err then
				return 1
			else
				cmd.Close()
			end
		else
			cmd.Close()
		end
	end
end
