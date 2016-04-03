-- exec wrappers

--Default read size
exec.READ_SIZE = 64

function exec.source(...) -- LTN12 compatible exec call! \o/
	local cmd = exec.exec(...)
	return function()
		local line, err = cmd.Read_Stdout(exec.READ_SIZE)
		if not err then
			return line
		end
	end
end

function exec.filter(...) -- returns two ltn12 functions, a sink and a source respectively.
	local cmd = exec.exec(...)
	local sink = function(chunk)
		if chunk then
			local err = cmd.Write_Stdin(tostring(chunk))
			if not err then
				return 1
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

function exec.sink(...)
	local cmd = exec.exec(...)
	return function(chunk)
		if chunk then
			local err = cmd.Write_Stdin(tostring(chunk))
			if not err then
				return 1
			end
		end
	end
end
