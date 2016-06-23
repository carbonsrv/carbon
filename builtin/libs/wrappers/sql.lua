-- SQL library
sql = {}

-- Helpers
local function do_ping(self)
	if not self.in_tranaction then
		local err = self.main_con.Ping()
		if err then
			error("SQL: Database connection ping failed: "..err, 0)
		end
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
		-- Basics
		main_con = db,
		con = db, -- switched when transaction
		in_transaction = false,
		close = function(self)
			local err = self.con.Close()
			if err then
				return false, err
			end
			return true, nil
		end,
		-- Transactions
		begin = function(self, fn)
			if self.in_tranaction then
				error("SQL: Already in transaction!", 0)
			end

			do_ping(self)

			self.in_transaction = true
			local tx, err = self.main_con.Begin()
			if err then
				return nil, err
			end
			self.con = tx

			if fn then
				local docommit = fn(self)
				if docommit == true then
					return self:commit()
				elseif docommit == false then
					return self:rollback()
				else
					self:rollback()
					error("SQL: begin function did not return true (commit) or false (rollback), rolling back unconditionally", 0)
				end
			end
		end,
		commit = function(self)
			if not self.in_transaction then
				error("SQL: not in transaction, can't commit", 0)
			end

			local e = self.con.Commit()
			self.con = self.main_con
			if e then
				return false, e
			end
			return true, nil
		end,
		rollback = function(self)
			if not self.in_transaction then
				error("SQL: not in transaction, can't rollback", 0)
			end

			local e = self.con.Rollback()
			self.con = self.main_con
			if e then
				return false, e
			end
			return true, nil
		end,
		-- Standard things.
		exec = function(self, stmt, ...)
			if not stmt then
				error("SQL: exec needs statemement!", 0)
			end

			do_ping(self)
			return self.con.Exec(stmt, ...)
		end,
		query = function(self, stmt, ...)
			if not stmt then
				error("SQL: query needs statement!", 0)
			end

			do_ping(self)
			local rows, err = self.con.Query(stmt, ...)
			if err then
				return nil, err
			end

			return convert_rows(rows)
		end,
		-- Prepared statements
		prepare = function(origself, stmt)
			if not stmt then
				error("SQL: prepare needs statement!", 0)
			end

			do_ping(origself)
			local pstmt, err = origself.con.Prepare(stmt)
			if err then
				return nil, err
			end

			return {
				["main_con"] = origself.con,
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
