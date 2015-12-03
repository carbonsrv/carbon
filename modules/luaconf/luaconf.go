package luaconf

import (
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/gin-gonic/gin"
	//"github.com/vifino/golua/lua"
	"github.com/carbonsrv/carbon/modules/middleware"
	"github.com/carbonsrv/carbon/modules/repl"
	"github.com/pmylund/go-cache"
	"github.com/vifino/luar"
)

// Configure the server based on a lua script.
func Configure(script string, args []string, cfe *cache.Cache, webroot string, useRecovery bool, useLogger bool, runrepl bool, finalizer func(srv *gin.Engine)) error {
	srv := gin.New()
	if useLogger {
		srv.Use(gin.Logger())
	}
	if useRecovery {
		srv.Use(gin.Recovery())
	}
	L := luar.Init()
	luar.Register(L, "carbon", luar.Map{ // srv and the state
		"srv": srv,
		"L":   L,
	})
	luar.Register(L, "var", luar.Map{ // Vars
		"root": webroot,
	})
	luar.Register(L, "", luar.Map{
		"args": args,
	})
	middleware.Bind(L)
	middleware.BindStatic(L, cfe)
	L.DoString(glue.MainGlue())
	L.DoString(glue.ConfGlue())
	go finalizer(srv)
	err := L.DoFile(script)
	if err != nil {
		return err
	} else {
		if runrepl {
			repl.Run(L)
		}
		c := make(chan bool)
		<-c
	}
	return nil
}

// Configure the server based on a lua string.
func Eval(script string, args []string, cfe *cache.Cache, webroot string, useRecovery bool, useLogger bool, runrepl bool, finalizer func(srv *gin.Engine)) error {
	srv := gin.New()
	if useLogger {
		srv.Use(gin.Logger())
	}
	if useRecovery {
		srv.Use(gin.Recovery())
	}
	L := luar.Init()
	luar.Register(L, "carbon", luar.Map{ // srv and the state
		"srv": srv,
		"L":   L,
	})
	luar.Register(L, "var", luar.Map{ // Vars
		"root": webroot,
	})
	luar.Register(L, "", luar.Map{
		"args": args,
	})
	middleware.Bind(L)
	middleware.BindStatic(L, cfe)
	L.DoString(glue.MainGlue())
	L.DoString(glue.ConfGlue())
	go finalizer(srv)
	err := L.DoString(script)
	if err != nil {
		return err
	} else {
		if runrepl {
			repl.Run(L)
		}
		c := make(chan bool)
		<-c
	}
	return nil
}

// Run REPL
func REPL(args []string, cfe *cache.Cache, webroot string, useRecovery bool, useLogger bool, finalizer func(srv *gin.Engine)) error {
	srv := gin.New()
	if useLogger {
		srv.Use(gin.Logger())
	}
	if useRecovery {
		srv.Use(gin.Recovery())
	}
	L := luar.Init()
	luar.Register(L, "carbon", luar.Map{ // srv and the state
		"srv": srv,
		"L":   L,
	})
	luar.Register(L, "var", luar.Map{ // Vars
		"root": webroot,
	})
	luar.Register(L, "", luar.Map{
		"args": args,
	})
	middleware.Bind(L)
	middleware.BindStatic(L, cfe)
	L.DoString(glue.MainGlue())
	L.DoString(glue.ConfGlue())
	go finalizer(srv)

	repl.Run(L)
	return nil
}
