package repl

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"strings"
)

// Dummy struct for the Run method
type Dummy struct {
	Name string
}

const (
	Lua_Prompt1 = "> "
	Lua_Prompt2 = ">> "
	Banner      = `
  ⡎⠑ | Carbon
  ⠣⠔ | (c) 2015 vifino
`
)

// modified https://github.com/vifino/luar/blob/master/examples/luar.go
// Run the repl
func Run(L *lua.State) {
	luar.Register(L, "", luar.Map{
		"__DUMMY__": &Dummy{"me"},
	})

	err := L.DoString(lua_code)
	if err != nil {
		fmt.Println("initial " + err.Error())
		return
	}
	// particularly the completion logic
	complete := luar.NewLuaObjectFromName(L, "carbon._lua_candidates")
	// this function returns a string slice of candidates
	str_slice := luar.Types([]string{})

	linenoise.SetCompletionHandler(func(in string) []string {
		val, err := complete.Callf(str_slice, in)
		if err != nil || len(val) == 1 && val[0] == nil {
			return []string{}
		}
		return val[0].([]string)
	})

	fmt.Println(Banner)
	L.DoString(`print(jit.version)`)
	fmt.Println("")

	prompt := Lua_Prompt1
	code := ""

	for {
		// ctrl-C/ctrl-D handling with latest update of go.linenoise
		str, err := linenoise.Line(prompt)
		if err != nil {
			return
		}
		if len(str) > 0 {
			linenoise.AddHistory(str)
			if str[0] == '=' || str[0] == '.' {
				exprs := str[1:]
				if str[0] == '=' {
					str = "print(serialize.simple(" + exprs + "))"
				} else {
					str = "print(" + exprs + ")"
				}
			}
			continuing := false
			code = code + str
			//fmt.Println("'"+code+"'")
			err := L.DoString(code)
			if err != nil {
				errs := err.Error()
				// fmt.Println("err",errs)
				// strip line nonsense if error occurred on prompt
				idx := strings.Index(errs, ": ")
				if idx > -1 && strings.HasPrefix(errs, "[string ") {
					errs = errs[idx+2:]
				}
				// need to keep collecting line?
				if strings.HasSuffix(errs, "near '<eof>'") {
					continuing = true
				} else {
					fmt.Println("Error: " + errs)
				}
			}
			if continuing { // prompt must reflect continuing state
				prompt = Lua_Prompt2
				code = code + "\n"
			} else {
				prompt = Lua_Prompt1
				code = ""
			}
		}
	}
}

// pretty-printing and code completion logic in Lua
const lua_code = `
serialize = require("serialize")

local tostring = tostring
local append = table.insert
local function quote (v)
	local t = type(v)
	if t == 'string' then
		return ('%q'):format(v)
	elseif t == 'function' then
		return '<fun>'
	elseif t == 'userdata' then
		return '<udata>'
	else
		return tostring(v)
	end
end
local dump
dump = function(t, options)
	options = options or { }
	local limit = options.limit or 1000
	local buff = {tables={}}
	if type(t) == 'table' then
			buff.tables[t] = true
	end
	local k, tbuff = 1, nil
	local function put(v)
		buff[k] = v
		k = k + 1
	end
	local function put_value(value)
		if type(value) ~= 'table' then
			put(quote(value))
			if limit and k > limit then
				buff[k] = "..."
				error("buffer overrun")
			end
		else
			if not buff.tables[value] then
				buff.tables[value] = true
				tbuff(value)
			else
				put("<cycle>")
			end
		end
		return put(',')
	end
	function tbuff(t)
		local mt
		if not (options.raw) then
			mt = getmetatable(t)
		end
		if mt and mt.__tostring then
			return put(quote(t))
		elseif type(t) ~= 'table' and not (mt and mt.__pairs) then
			return put(quote(t))
		else
			put('{')
			local mt_pairs, indices = mt and mt.__pairs
			if not mt_pairs and #t > 0 then
				indices = {}
				for i = 1, #t do
					indices[i] = true
				end
			end
			for key, value in pairs(t) do
				local _continue_0 = false
				repeat
					if indices and indices[key] then
						_continue_0 = true
						break
					end
					if type(key) ~= 'string' then
						key = '[' .. tostring(key) .. ']'
					elseif key:match('%s') then
						key = quote(key)
					end
					put(key .. '=')
					put_value(value)
					_continue_0 = true
				until true
				if not _continue_0 then
					break
				end
			end
			if indices then
				local _list_0 = t
				for _index_0 = 1, #_list_0 do
					local v = _list_0[_index_0]
					put_value(v)
				end
			end
			if buff[k - 1] == "," then
				k = k - 1
			end
			return put('}')
		end
	end
	tbuff(t)
	--pcall(tbuff, t)
	return table.concat(buff)
end

local append = table.insert
local function is_pair_iterable(t)
		local mt = getmetatable(t)
		return type(t) == 'table' or (mt and mt.__pairs)
end

function carbon._lua_candidates(line)
	-- identify the expression!
	local res = {}
	local i1,i2 = line:find('[.:%w_]+$')
	if not i1 then return res end
	local front,partial = line:sub(1,i1-1), line:sub(i1)
	local prefix, last = partial:match '(.-)([^.:]*)$'
	local t, all = _G
	if #prefix > 0 then
		local P = prefix:sub(1,-2)
		all = last == ''
		for w in P:gmatch '[^.:]+' do
			t = t[w]
			if not t then
				return res
			end
		end
	end
	prefix = front .. prefix
	local function append_candidates(t)
		for k,v in pairs(t) do
			if all or k:sub(1,#last) == last then
				append(res,prefix..k)
			end
		end
	end
	local mt = getmetatable(t)
	if is_pair_iterable(t) then
		append_candidates(t)
	end
	if mt and is_pair_iterable(mt.__index) then
		append_candidates(mt.__index)
	end
	return res
end

--// override struct __pairs for code completion
local function sdump(st)
		local t = luar.type(st)
		local val = luar.value(st)
		local nm = t.NumMethod()
		local mt = t --// type to dispatch methods on ptr receiver
		if t.Kind() == 22 then --// pointer!
				t = t.Elem()
				val = val.Elem()
		end
		local n = t.NumField()
		local cc = {}
		for i = 1,n do
				local f,v = t.Field(i-1)
				if f.PkgPath == "" then --// only public fields!
						v = val.Field(i-1)
						cc[f.Name] = v.Interface()
				end
		end
		--// then public methods...
		for i = 1,nm do
				local m = mt.Method(i-1)
				if m.PkgPath == "" then --// again, only public
						cc[m.Name] = true
				end
		end
		return cc
end

mt = getmetatable(__DUMMY__)
mt.__pairs = function(st)
		local cc = sdump(st)
		return pairs(cc)
end
`
