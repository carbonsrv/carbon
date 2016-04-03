-- termbox
-- Small helpers
function termbox.print(x, y, msg, fg, bg)
	local fg = fg or termbox.ColorWhite
	local bg = bg or termbox.ColorDefault
	termbox.TBPrint(x, y, fg, bg, msg)
end
