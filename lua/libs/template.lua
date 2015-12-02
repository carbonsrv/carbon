-- Template rendering

local template = {}

local escapist = require("escapist")
local prettify = require("serialize").simple

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

function template.render(source, env)
	local env = env or _G

	if source then
		local output = tostring(source):gsub("<%%%%&(.-)%%%%>(.-)<%%&%%>", function(content, src) -- Have an intermediary template, render that, escaping html and providing <%%& codehere %%> <%.%> <%&%>
			local res = ""
			local f = eval(content, env)
			local result = {pcall(f)}
			local suc = table.remove(result, 1)
			if suc then
				local iter = table.orderedPairs
				local t
				if type(result[1]) == "table" then
					t = result[1]

					if t[1] then -- If it seems to be n based..
						iter = ipairs
					end
				elseif type(result[1]) == "function" and type(result[2]) == "table" then -- if it is the result of an iterator
					iter = table.unpack
					t = result
				else
					t = result
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
				return "", result[1]
			end
		end):gsub("<%%&(.-)%%>(.-)<%%&%%>", function(content, src) -- Have an intermediary template, render that, escaping html and providing <%& codehere %> <%.%> <%&%>
			local res = ""
			local f = eval(content, env)
			local result = {pcall(f)}
			local suc = table.remove(result, 1)
			if suc then
				local iter = table.orderedPairs
				local t
				if type(result[1]) == "table" then
					t = result[1]

					if t[1] then -- If it seems to be n based..
						iter = ipairs
					end
				elseif type(result[1]) == "function" and type(result[2]) == "table" then -- if it is the result of an iterator
					iter = table.unpack
					t = result
				else
					t = result
					iter = ipairs
				end

				for k, v in iter(t) do
					local eenv = env
					eenv["k"] = k
					eenv["v"] = v
					local k2 = escapist.escape.html(prettify(k))
					local v2 = escapist.escape.html(prettify(v))
					local r, err = template.render(src:gsub("<%%%.%%>", v2):gsub("<%%k%%>", k2):gsub("<%%v%%>", v2), eenv)
					if err then
						return "", err
					end
					res = res .. r
				end
				return res
			else
				return "", result[1]
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
