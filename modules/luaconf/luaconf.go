package luaconf

import (
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/gin-gonic/gin"
	//"github.com/vifino/golua/lua"
	"github.com/carbonsrv/carbon/modules/middleware"
	"github.com/carbonsrv/carbon/modules/repl"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/pmylund/go-cache"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
)

var checker_code = `
srv.Finish_VHOSTS()

if srv.used == false then
	os.exit(0)
end
`

// Setup basic state
func Setup(args []string, cfe *cache.Cache, webroot string, useRecovery bool, useLogger bool, runrepl bool, finalizer func(srv *gin.Engine), bindhook func(L *lua.State)) *lua.State {
	srv := gin.New()
	if useLogger {
		srv.Use(gin.Logger())
	}
	if useRecovery {
		srv.Use(gin.Recovery())
	}
	L := luar.Init()

	var didnt_run_yet = true
	luar.Register(L, "carbon", luar.Map{ // srv and the state
		"srv": srv,
		"L":   L,
		"launch_server": func() {
			if didnt_run_yet {
				scheduler.Add(func() {
					finalizer(srv)
				})
				didnt_run_yet = false
			}
		},
	})
	luar.Register(L, "", luar.Map{
		"arg": args,
	})
	middleware.Bind(L, webroot)
	middleware.BindStatic(L, cfe)
	bindhook(L)
	L.DoString(glue.MainGlue())
	L.DoString(glue.ConfGlue())
	return L
}

// Configure the server based on a lua script.
func Configure(script string, args []string, cfe *cache.Cache, webroot string, useRecovery bool, useLogger bool, runrepl bool, finalizer func(srv *gin.Engine), bindhook func(L *lua.State)) error {
	L := Setup(args, cfe, webroot, useRecovery, useLogger, runrepl, finalizer, bindhook)

	err := L.DoFile(script)
	if err == nil {
		if runrepl {
			repl.Run(L)
		}
		L.DoString(checker_code)
		c := make(chan bool)
		<-c
	}
	L.Close()
	return err
}

// Eval lua string to Configure the server
func Eval(script string, args []string, cfe *cache.Cache, webroot string, useRecovery bool, useLogger bool, runrepl bool, finalizer func(srv *gin.Engine), bindhook func(L *lua.State)) error {
	L := Setup(args, cfe, webroot, useRecovery, useLogger, runrepl, finalizer, bindhook)

	err := L.DoString(script)
	if err == nil {
		if runrepl {
			repl.Run(L)
		}
		L.DoString(checker_code)
		c := make(chan bool)
		<-c
	}
	L.Close()
	return err
}

// REPL runs a lua repl
func REPL(args []string, cfe *cache.Cache, webroot string, useRecovery bool, useLogger bool, finalizer func(srv *gin.Engine), bindhook func(L *lua.State)) error {
	L := Setup(args, cfe, webroot, useRecovery, useLogger, true, finalizer, bindhook)

	repl.Run(L)
	L.DoString(checker_code)
	L.Close()
	return nil
}
