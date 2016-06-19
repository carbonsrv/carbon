package repl

import (
	"fmt"
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
)

// Dummy struct for the Run method
type Dummy struct {
	Name string
}

// heavily inspired by https://github.com/vifino/luar/blob/master/examples/luar.go
// Run the repl
func Run(L *lua.State) {
	luar.Register(L, "", luar.Map{
		"carbon.__DUMMY__": &Dummy{"me"},
	})

	err := L.DoString(glue.GetGlue("REPL.lua"))
	if err != nil {
		fmt.Println("initial " + err.Error())
		return
	}
}
