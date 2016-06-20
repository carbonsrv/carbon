-- Print golang blog posts
local yql = sql.open("yql")

local rows = assert(yql:query("select * from atom where url = ?", "http://blog.golang.org/feeds/posts/default?alt=rss"))

for row=1, rows.n do
	print("Post no. "..tostring(row)..":")
	local row = rows[row].results
	print(row.title..": https:"..row.link.href)
	print()
end
