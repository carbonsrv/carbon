package glue

// MainGlue returns the main glue
func MainGlue() string {
	data, err := Asset("MainGlue.lua")
	if err != nil {
		panic(err)
	}
	return string(data)
}

// RouteGlue returns the glue for routes
func RouteGlue() string {
	data, err := Asset("RouteGlue.lua")
	if err != nil {
		panic(err)
	}
	return string(data)
}

// ConfGlue returns the glue for the config state
func ConfGlue() string {
	data, err := Asset("ConfGlue.lua")
	if err != nil {
		panic(err)
	}
	return string(data)
}

// GetGlue returns the builtin assets compiled in.
func GetGlue(asset string) string {
	data, err := Asset(asset)
	if err != nil {
		return ""
	}
	return string(data)
}
