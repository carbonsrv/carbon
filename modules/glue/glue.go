package glue

func Glue() string {
	data, _ := Asset("MainGlue.lua")
	return string(data)
}
