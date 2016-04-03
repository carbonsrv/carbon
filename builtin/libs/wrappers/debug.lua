-- debug enhancements
function debug.getallupvals(f)
	local i = 1
	local r = {}
	while true do
		local n, v = debug.getupvalue(f, i)
		if not n then
			if r == {} then
				return nil
			end
			return r
		end
		if n ~= "_ENV" then
			r[i] = {name=n, value=v}
		end
		i = i + 1
	end
end

function debug.setallupvals(f, vals)
	for i, pair in pairs(vals) do
		debug.setupvalue(f, i, pair["value"])
	end
end
