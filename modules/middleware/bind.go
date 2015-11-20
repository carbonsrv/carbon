package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/fzzy/radix/redis"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/carbonsrv/carbon/modules/helpers"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/carbonsrv/carbon/modules/static"
	"github.com/vifino/contrib/gzip"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"net"
	"time"
)

func Bind(L *lua.State) {
	BindCarbon(L)
	BindMiddleware(L)
	BindRedis(L)
	BindKVStore(L)
	BindPhysFS(L)
	BindThread(L)
	BindOther(L)
	BindNet(L)
	BindConversions(L)
	BindComs(L)
}

func BindCarbon(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // Carbon specific API
		"glue": glue.GetGlue,
	})
}

func BindEngine(L *lua.State) {
	luar.Register(L, "engine", luar.Map{
		"new": gin.New,
	})
}

func BindMiddleware(L *lua.State) {
	luar.Register(L, "mw", luar.Map{
		"Lua": Lua,
		"ExtRoute": (func(plan map[string]interface{}) func(*gin.Context) {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(func(*gin.Context))
			}
			return ExtRoute(newplan)
		}),
		"VHOST": (func(plan map[string]interface{}) func(*gin.Context) {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(func(*gin.Context))
			}
			return VHOST(newplan)
		}),
		"Logger":   gin.Logger,
		"Recovery": gin.Recovery,
		"GZip": func() func(*gin.Context) {
			return gzip.Gzip(gzip.DefaultCompression)
		},
		"DLR_NS":    DLR_NS,
		"DLR_RUS":   DLR_RUS,
		"DLRWS_RUS": DLRWS_RUS,
		"Echo":      EchoHTML,
		"EchoText":  Echo,
	})
	L.DoString(glue.RouteGlue())
}

func BindPhysFS(L *lua.State) {
	luar.Register(L, "fs", luar.Map{ // PhysFS
		"mount":       physfs.Mount,
		"exits":       physfs.Exists,
		"getFS":       physfs.FileSystem,
		"mkdir":       physfs.Mkdir,
		"umount":      physfs.RemoveFromSearchPath,
		"delete":      physfs.Delete,
		"setWriteDir": physfs.SetWriteDir,
		"getWriteDir": physfs.GetWriteDir,
	})
}

func BindIOEnhancements(L *lua.State) {
	luar.Register(L, "io", luar.Map{ // Small enhancements to the io stuff.
		"list": (func(path string) ([]string, error) {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return make([]string), err
			} else {
				list := make([]string, len(files))
				for i := range files {
					list[i] = files[i].Name()
				}
				return list, nil
			}
		}),
	})
}

func BindRedis(L *lua.State) {
	luar.Register(L, "redis", luar.Map{
		"connectTimeout": (func(host string, timeout int) (*redis.Client, error) {
			return redis.DialTimeout("tcp", host, time.Duration(timeout)*time.Second)
		}),
		"connect": (func(host string) (*redis.Client, error) {
			return redis.Dial("tcp", host)
		}),
	})
}

func BindKVStore(L *lua.State) { // Thread safe Key Value Store that doesn't persist.
	luar.Register(L, "kvstore", luar.Map{
		"set": (func(k string, v interface{}){
			kvstore.Set(k, v, -1)
		}),
		"del": (func(k string){
			kvstore.Delete(k)
		}),
		"get": (func(k string) interface{}{
			res, found := kvstore.Get(k)
			if found {
				return res
			} else {
				return nil
			}
		}),
		"inc": (func(k string, n int64) error{
			return kvstore.Increment(k, n)
		}),
		"dec": (func(k string, n int64) error{
			return kvstore.Decrement(k, n)
		}),
	})
}

func BindThread(L *lua.State) {
	luar.Register(L, "thread", luar.Map{
		"_spawn": (func(bcode string, dobind bool, vals map[string]interface{}) error {
			L := luar.Init()
			Bind(L)
			err := L.DoString(glue.MainGlue())
			if err != nil {
				panic(err)
			}

			if dobind {
				luar.Register(L, "", vals)
			}

			if L.LoadBuffer(bcode, len(bcode), "thread") != 0 {
				return errors.New(L.ToString(-1))
			}

			scheduler.Add(func() {
				if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
					// Silently error because reasons. ._.
				}
			})
			return nil
		}),
	})
}

func BindOther(L *lua.State) {
	luar.Register(L, "", luar.Map{
		"unixtime": (func() int {
			return int(time.Now().UnixNano())
		}),
		"_syntaxhlfunc": helpers.SyntaxHL,
	})
}

func BindComs(L *lua.State) {
	luar.Register(L, "com", luar.Map{
		"create": (func() chan interface{} {
			return make(chan interface{})
		}),
		"createBuffered": (func(buffer int) chan interface{} {
			return make(chan interface{}, buffer)
		}),
		"receive": (func(c chan interface{}) interface{} {
			return <-c
		}),
		"send": (func(c chan interface{}, val interface{}) {
			c <- val
		}),
	})
}

func BindNet(L *lua.State) {
	luar.Register(L, "net", luar.Map{
		"dial": net.Dial,
		"write": (func(con net.Conn, str string) {
			fmt.Fprintf(con, str)
		}),
		"readline": (func(con net.Conn) (string, error) {
			return bufio.NewReader(con).ReadString('\n')
		}),
	})
}

func BindConversions(L *lua.State) {
	luar.Register(L, "convert", luar.Map{
		"stringtocharslice": (func(x string) []byte {
			return []byte(x)
		}),
		"charslicetostring": (func(x []byte) string {
			return string(x)
		}),
	})
}

func BindContext(L *lua.State, context *gin.Context) {
	luar.Register(L, "", luar.Map{
		"context":    context,
		"req":        context.Request,
		"_paramfunc": context.Param,
		"_formfunc":  context.PostForm,
		"_queryfunc": context.Query,
	})
}
func BindStatic(L *lua.State, cfe *cache.Cache) {
	luar.Register(L, "static", luar.Map{
		"serve": (func(prefix string) func(*gin.Context) {
			return staticServe.ServeCached(prefix, staticServe.PhysFS("", true, true), cfe)
		}),
	})
}
