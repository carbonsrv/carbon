local db = sql.open("sqlite3", ":memory:")
print(db:exec("CREATE TABLE users (name text, age int)"))
print(db:exec("INSERT INTO users (name, age) VALUES ($1, $2)", "vifino", 16))
print(db:exec("INSERT INTO users (name, age) VALUES ($1, $2)", "Bob", 21))

local res = assert(db:query("SELECT * FROM users"))
for rowno=1, res.n
 do
	print("Showing row no. "..tostring(rowno))
	for k, v in pairs(res[rowno]) do
		print(k, v)
	end
end
