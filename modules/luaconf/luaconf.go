package luaconf

import (
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/gin-gonic/gin"
	//"github.com/vifino/golua/lua"
	"github.com/carbonsrv/carbon/modules/middleware"
	"github.com/pmylund/go-cache"
	"github.com/vifino/luar"
	"net/http"
)

// Configure the server based on a lua script.
func Configure(script string, cfe *cache.Cache, webroot string) (http.Handler, error) {
	srv := gin.New()
	L := luar.Init()
	luar.Register(L, "carbon", luar.Map{ // srv and the state
		"srv": srv,
		"L":   L,
	})
	luar.Register(L, "var", luar.Map{ // Vars
		"root": webroot,
	})
	middleware.Bind(L)
	middleware.BindStatic(L, cfe)
	L.DoString(glue.MainGlue())
	L.DoString(glue.ConfGlue())
	return srv, L.DoFile(script)
}
