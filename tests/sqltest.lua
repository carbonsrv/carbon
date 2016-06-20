local db = sql.open("ql-mem", "test.qldb") -- memory backed ql database

assert(db:begin(function() -- ql needs everything to be in a transaction
	assert(db:exec("CREATE TABLE users (name string, age float64)")) -- Create a new table called users with the rows name and age
	assert(db:exec("INSERT INTO users (name, age) VALUES ($1, $2)", "bauen1", 15))
	assert(db:exec("INSERT INTO users (name, age) VALUES ($1, $2)", "bob", 21))
	assert(db:exec("INSERT INTO users (name, age) VALUES ($1, $2)", "vifino", 16))
	return true -- commit
end))

print("Listing all users:")
local rows = assert(db:query("SELECT * FROM users"))
for i=1, rows.n do
  print("Showing row number " .. tostring(i))
  for k,v in pairs(rows[i]) do
    print("> " .. tostring(k), tostring(v))
  end
end

print()
local min_ages = {8, 18}
local statement = db:prepare("SELECT * FROM users WHERE age>=$1")
for _,min_age in pairs(min_ages) do
  print("Users with a minimum age of " .. tostring(min_age))
  local rows = statement:query(min_age)
  for i=1, rows.n do
      local row = rows[i]
      print("> " .. row.name..": ".. tostring(row.age))
  end
end
