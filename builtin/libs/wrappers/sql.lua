-- SQL library
sql = {}

-- Helpers
local function do_ping(con)
	local err = con.Ping()
	if err then
		error("SQL: Database connection ping failed.", 0)
	end
end
local function convert_rows(rows)
	if not rows then
		error("SQL: Can't convert nil to row!")
	end

	local rows, err = carbon._sql_rows(rows)
	if err then
		error("SQL: Row conversion failed: "..err) 
	end

	local final = {}
	final.n = rows.Len
	for i=1, final.n do
		final[i] = luar.map2table(carbon._sql_getrow(rows, i))
	end

	return final
end

-- API
function sql.drivers() -- returns list of drivers available
	return luar.slice2table(carbon._sql_drivers())
end

function sql.open(driver, src) -- generates a database wrapper
	if not driver or not src then
		error("SQL: open needs driver and src", 0)
	end

	local db, err = carbon._sql_open(driver, src)
	if err then
		return nil, err
	end
	
	return {
		con = db,
		close = function()
			local err = self.con.Close()
			if err then
				return false, err
			end
			return true, nil
		end,
		exec = function(self, stmt, ...)
			if not stmt then
				error("SQL: exec needs statemement!", 0)
			end

			do_ping(self.con)
			return self.con.Exec(stmt, ...)
		end,
		query = function(self, stmt, ...)
			if not stmt then
				error("SQL: query needs statement!", 0)
			end

			do_ping(self.con)
			local rows, err = self.con.Query(stmt, ...)
			if err then
				return nil, err
			end

			return convert_rows(rows)
		end,
		prepare = function(stmt)
			if not stmt then
				error("SQL: prepare needs statement!", 0)
			end

			do_ping(self.con)
			local pstmt, err = self.con.Prepare(stmt)
			if err then
				return nil, err
			end

			return function(...)
				local rows, err = pstmt.Query(...)
				if err then
					return nil, err
				end

				return convert_rows(rows)
			end
		end,
	}
end
