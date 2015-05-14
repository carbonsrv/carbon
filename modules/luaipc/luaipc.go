package luaipc

import (
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
)

type IPC chan interface{}

func BindValues(L *lua.State, vals map[string]interface{}) {
	luar.Register(L, "", vals)
}

func New() IPC {
	return make(IPC)
}

func Send(c IPC, v interface{}) {
	c <- v
}

func Receive(c IPC) interface{} {
	return <-c
}
