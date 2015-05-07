package glue

func MainGlue() string {
	data, err := Asset("MainGlue.lua")
	if err != nil {
		panic(err)
	}
	return string(data)
}
func RouteGlue() string {
	data, err := Asset("RouteGlue.lua")
	if err != nil {
		panic(err)
	}
	return string(data)
}
