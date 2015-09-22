package middleware

import (
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/fzzy/radix/redis"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/vifino/carbon/modules/glue"
	"github.com/vifino/carbon/modules/helpers"
	"github.com/vifino/carbon/modules/static"
	"github.com/vifino/contrib/gzip"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"time"
	"net"
	"fmt"
	"bufio"
)

func Bind(L *lua.State) {
	BindMiddleware(L)
	BindRedis(L)
	BindPhysFS(L)
	BindOther(L)
	BindNet(L)
	BindConversions(L)
	BindComs(L)
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
