-- Template rendering

local template = {}

local escapist = require("escapist")

local function eval(str, env)
	local f, err = loadstring("return ".. str, "template")
	if err then
		f, err = loadstring(str, "template")
	end
	if err then
		error(err, nil)
	end

	return setfenv(f, env)
end

function prettify(...) -- Almost JSON(tm)
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
						prettifyrecursive(v, tab + #k + 2, true)
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


function template.render(source, env)
	local env = env or _G

	if source then
		local output = tostring(source):gsub("<%%%%&(.-)%%%%>(.-)<%%&%%>", function(content, src) -- Have an intermediary template, render that, escaping html and providing <%%& codehere %%> <%.%> <%&%>
			local res = ""
			local f = eval(content, env)
			local result = {pcall(f)}
			local suc = table.remove(result, 1)
			if suc then
				local t
				if type(result[1]) == "table" then
					t = result[1]
				else
					t = result
				end

				local iter = pairs
				if t[1] then -- If it seems to be n based..
					iter = ipairs
				end

				for k, v in iter(t) do
					local eenv = env
					eenv["k"] = k
					eenv["v"] = v
					local k2 = prettify(k)
					local v2 = prettify(v)
					local r, err = template.render(src:gsub("<%%%.%%>", v2):gsub("<%%k%%>", k2):gsub("<%%v%%>", v2), eenv)
					if err then
						return "", err
					end
					res = res .. r
				end
				return res
			else
				return "", result[2]
			end
		end):gsub("<%%&(.-)%%>(.-)<%%&%%>", function(content, src) -- Have an intermediary template, render that, escaping html and providing <%& codehere %> <%.%> <%&%>
			local res = ""
			local f = eval(content, env)
			local result = {pcall(f)}
			local suc = table.remove(result, 1)
			if suc then
				local t
				if type(result[1]) == "table" then
					t = result[1]
				else
					t = result
				end

				local iter = pairs
				if t[1] then -- If it seems to be n based..
					iter = ipairs
				end

				for k, v in iter(t) do
					local eenv = env
					eenv["k"] = k
					eenv["v"] = v
					local k2 = escapist.escape.html(prettify(k))
					local v2 = escapist.escape.html(prettify(v))
					res = res .. template.render(src:gsub("<%%%.%%>", v2):gsub("<%%k%%>", k2):gsub("<%%v%%>", v2), eenv)
				end
				return res
			else
				return "", result[2]
			end
		end):gsub("<%%%%=(.-)%%%%>", function(content) -- Prettify output, provides <%%= codehere %%>
			local res = ""
			local f = eval(content, env)
			local suc, result = pcall(f)
			if suc then
				return prettify(result)
			else
				return "", result
			end
		end):gsub("<%%=(.-)%%>", function(content) -- Prettify output and escape, provides <%= codehere %>
			local res = ""
			local f = eval(content, env)
			local suc, result = pcall(f)
			if suc then
				return escapist.escape.html(prettify(result)):gsub("\n", "<br />")
			else
				return "", result
			end
		end)
		return output
	else
		error("No Template Source given!")
	end
end

return template
