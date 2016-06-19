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

function sql.open(driver, dsn) -- generates a database wrapper
	if not driver then
		error("SQL: open needs at least driver", 0)
	end

	dsn = dsn or ""

	local db, err = carbon._sql_open(driver, dsn)
	if err then
		return nil, err
	end
	
	return {
		con = db,
		close = function(self)
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
		prepare = function(origself, stmt)
			if not stmt then
				error("SQL: prepare needs statement!", 0)
			end

			do_ping(origself.con)
			local pstmt, err = origself.con.Prepare(stmt)
			if err then
				return nil, err
			end

			return {
				["con"] = origself.con,
				["pstmt"] = pstmt,
				query = function(self, ...)
					local rows, err = self.pstmt.Query(...)
					if err then
						return nil, err
					end

					return convert_rows(rows)
				end,
				exec = function(self, ...)
					return self.pstmt.Exec(...)
				end,
				close = function(self)
					local err = self.pstmt.Close()
					if err then
						return false, err
					end
					return true, nil
				end,
			}
		end,
	}
end
