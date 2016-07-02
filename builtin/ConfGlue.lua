-- Config Glue
-- Runs once, basically.

-- Arguments
arg = arg and luar.slice2table(arg)
args = arg

-- Server wrappers
srv = srv or require("wrappers.srv")

require("wrappers.mw")

-- Initialization and things.
require("wrappers.init")
