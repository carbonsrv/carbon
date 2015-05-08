package middleware

import (
	"../glue"
	"../static"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
)

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
		"Logger":   gin.Logger,
		"Recovery": gin.Recovery,
		"GZip": func() func(*gin.Context) {
			return gzip.Gzip(gzip.DefaultCompression)
		},
		"DLR_NS":   DLR_NS,
		"DLR_RUS":  DLR_RUS,
		"Echo":     Echo,
		"EchoHTML": EchoHTML,
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
func BindStatic(L *lua.State, cfe *cache.Cache) {
	luar.Register(L, "static", luar.Map{
		"serve": (func(prefix string) func(*gin.Context) {
			return staticServe.ServeCached(prefix, staticServe.PhysFS("", true, true), cfe)
		}),
	})
}
