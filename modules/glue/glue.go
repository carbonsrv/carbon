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
func ConfGlue() string {
	data, err := Asset("ConfGlue.lua")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func GetGlue(asset string) string {
	data, err := Asset(asset)
	if err != nil {
		return nil
	}
	return string(data)
}
