package luaconf

import (
	"../routes"
	"../static"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	//"github.com/vifino/golua/lua"
	"github.com/pmylund/go-cache"
	"github.com/vifino/luar"
)

// Configure the server based on a lua script.
func Configure(script string, cfe *cache.Cache, webroot string) (*gin.Engine, error) {
	srv := gin.New()
	L := luar.Init()
	luar.Register(L, "", luar.Map{ // Global
		"srv": srv,
		"L":   L,
	})
	luar.Register(L, "var", luar.Map{ // Vars
		"root": webroot,
	})
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
	luar.Register(L, "mw", luar.Map{
		"Lua": routes.Lua,
		"ExtRoute": (func(plan map[string]interface{}) func(*gin.Context) {
			newplan := make(routes.Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(func(*gin.Context))
			}
			return routes.ExtRoute(newplan)
		}),
		"Logger":   gin.Logger,
		"Recovery": gin.Recovery,
		"GZip": func() func(*gin.Context) {
			return gzip.Gzip(gzip.DefaultCompression)
		},
	})
	luar.Register(L, "static", luar.Map{
		"serve": (func(prefix string) func(*gin.Context) {
			return staticServe.ServeCached(prefix, staticServe.PhysFS("", true, true), cfe)
		}),
	})
	return srv, L.DoFile(script)
}
