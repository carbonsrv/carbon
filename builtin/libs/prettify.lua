-- Prettify values for human inspection.

return function(...) -- Almost JSON(tm)
	local out = ""
	local objs = {...}
	local function prettifyrecursive(obj, tab, skipfirstindent)
		local tab = tab or 0
		local tabify = function(s) for i=1,tab do out = out .. ' ' end out = out .. s end
		local line = function(s) tabify(s .. '\n') end
		local mt = getmetatable(obj)
		if mt and mt.__tostring then
			out = out .. tostring(obj)
		else
			if not skipfirstindent then
				line('{')
			else
				out = out .. '{\n'
			end
			tab = tab+2
			for k,v in pairs(obj) do
				if type(v) == 'table' then
					if tab > 16 or next(v) == nil then
						line(k .. ': ' .. tostring(v))
					else
						tabify(k .. ': ')
						prettifyrecursive(v, tab + 2, true)
					end
				else
					line(k .. ': ' .. tostring(v))
				end
			end
			tab = tab-2
			line('}')
		end
	end
	for i = 1,select('#',...) do
		local obj = select(i,...)
		if type(obj) ~= 'table' then
			if type(obj) == 'userdata' or type(obj) == 'cdata' then
				out = out .. tostring(obj) .. '\n'
			else
				out = out .. tostring(obj) .. '\n'
			end
		else
			prettifyrecursive(obj)
		end
	end
	return out:sub(1, -2)
end
