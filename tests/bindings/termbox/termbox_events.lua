print("starting")
local err = termbox.Init()
if err then
	error(err)
end

termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
termbox.print(1, 1, "Hello world!", termbox.ColorRed)
termbox.Flush()
--print("test")
os.sleep(2)
while true do
	local e, err = termbox.PollEvent()
	if err then
		error(err)
	end
	
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.print(1, 1, "Hello world.")
	local row = 1
	for k, v in pairs(e) do
		termbox.print(1, 1+row, k..": "..tostring(v))
		row = row + 1
	end
	termbox.Flush()
end
