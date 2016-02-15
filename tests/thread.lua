thread = require("thread")
local a = "test"
thread.spawn(function()
	print(a)
end)
os.sleep(10)
