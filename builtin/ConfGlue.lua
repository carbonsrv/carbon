-- Config Glue

arg = arg and luar.slice2table(arg)
args = arg

srv = srv or require("wrappers.srv")

require("wrappers.mw")

require("wrappers.init")
