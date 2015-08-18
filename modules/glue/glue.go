package glue

func MainGlue() string {
	data, err := Asset("MainGlue.lua")
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
