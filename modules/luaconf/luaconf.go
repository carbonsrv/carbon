package luaconf

import (
	"../middleware"
	"github.com/gin-gonic/gin"
	//"github.com/vifino/golua/lua"
	"../glue"
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
	middleware.BindMiddleware(L)
	middleware.BindRedis(L)
	middleware.BindPhysFS(L)
	middleware.BindOther(L)
	middleware.BindStatic(L, cfe)
	L.DoString(glue.MainGlue())
	L.DoString(glue.RouteGlue())
	return srv, L.DoFile(script)
}
