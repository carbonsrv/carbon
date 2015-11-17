package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"

	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/carbonsrv/carbon/modules/helpers"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/carbonsrv/carbon/modules/static"
	"github.com/vifino/contrib/gzip"

	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/fzzy/radix/redis"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func Bind(L *lua.State) {
	BindCarbon(L)
	BindMiddleware(L)
	BindRedis(L)
	BindDB(L)
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

func BindDB(L *lua.State) {
	luar.Register(L, "db", luar.Map{
		"open": sql.Open,
		"rows": (func(rows *sql.Rows) (map[int]map[string]interface{}, error) {
			var res map[int]map[int]interface{}
			rowno := 1
			for rows.Next() {
				var elem1 interface{}
				var elem2 interface{}
				var elem3 interface{}
				var elem4 interface{}
				var elem5 interface{}
				var elem6 interface{}
				var elem7 interface{}
				var elem8 interface{}
				var elem9 interface{}
				var elem10 interface{}

				err := rows.Scan(&elem1, &elem2, &elem3, &elem4, &elem5, &elem6, &elem7, &elem8, &elem9, &elem10)
				if err != nil {
					return res, err
				}
				rowtmp := make(map[int]interface{})
				rowtmp[1] = elem1
				rowtmp[2] = elem2
				rowtmp[3] = elem3
				rowtmp[4] = elem4
				rowtmp[5] = elem5
				rowtmp[6] = elem6
				rowtmp[7] = elem7
				rowtmp[8] = elem8
				rowtmp[9] = elem9
				rowtmp[10] = elem10

				res[rowno] = rowtmp
				rowno += 1
			}
			err = rows.Err()
			if err != nil {
				return res, err
			} else {
				return res, nil
			}
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
