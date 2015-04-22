package main

//go:generate go-bindata -o glue/generated_lua.go -pkg=glue -prefix "./lua" ./lua

import (
	"./glue"
	"fmt"
	//"github.com/gin-gonic/contrib/gzip"
	"./modules/static"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"io/ioutil"
	//"net/http"
	//"runtime"
	"./scheduler"
	"errors"
	"github.com/namsral/flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Preloader/Starter
var jobs *int
var preloaded chan *lua.State

func preloader() {
	preloaded = make(chan *lua.State, *jobs)
	for {
		state := luar.Init()
		err := state.DoString(glue.Glue())
		if err != nil {
			fmt.Println(err)
		}
		preloaded <- state
	}
}
func getInstance() *lua.State {
	return <-preloaded
}

// Cache
var cfe *cache.Cache
var cbc *cache.Cache

func cacheRead(c *cache.Cache, file string) (string, error) {
	res := ""
	data_tmp, found := c.Get(file)
	if found == false {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		res = string(data)
		c.Set(file, res, cache.DefaultExpiration)
	} else {
		debug("Using cache for %s" + file)
		res = data_tmp.(string)
	}
	return res, nil
}
func cacheDump(L *lua.State, file string) (string, error, bool) {
	data_tmp, found := cbc.Get(file)
	if found == false {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err, false
		}
		if L.LoadString(string(data)) != 0 {
			return "", errors.New(L.ToString(-1)), true
		}
		res := L.FDump()
		L.Pop(1)
		cbc.Set(file, res, cache.DefaultExpiration)
		return res, nil, false
	} else {
		debug("Using Bytecode-cache for " + file)
		return data_tmp.(string), nil, false
	}
}
func cacheFileExists(file string) bool {
	data_tmp, found := cfe.Get(file)
	if found == false {
		exists := fileExists(file)
		cfe.Set(file, exists, cache.DefaultExpiration)
		return exists
	} else {
		return data_tmp.(bool)
	}
}

// Helper functions
func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

// Logging
var doLog bool = false

func debug(str string) {
	if doLog {
		log.Print(str)
	}
}

// Server
func new_server() *gin.Engine {
	r := gin.New()
	return r
}
func bootstrap(srv *gin.Engine, dir string) *gin.Engine {
	go preloader()     // Run the instance starter.
	go scheduler.Run() // Run the scheduler.
	srv.GET(`/:file`, logic_switcheroo(dir))
	//srv.Use(martini.Static(dir))
	return srv
}

// Routes
func logic_switcheroo(dir string) func(*gin.Context) {
	st := staticServe.ServeCached("/", staticServe.LocalFile(dir, true))
	lr := luaroute(dir)
	return func(context *gin.Context) {
		file := context.Params.ByName("file")
		fe := cacheFileExists(file)
		if fe == true {
			if strings.HasSuffix(file, ".lua") {
				lr(context)
			} else {
				st(context)
			}
		} else {
			context.String(404, "404 page not found")
		}
	}
}

func luaroute(dir string) func(*gin.Context) {
	LDumper := luar.Init()
	return func(context *gin.Context) {
		L := getInstance()
		defer scheduler.Add(func() {
			L.Close()
		})
		file := dir + context.Request.URL.Path
		luar.Register(L, "", luar.Map{
			"context": context,
			"req":     context.Request,
		})
		code, err, lerr := cacheDump(LDumper, file)
		if err != nil {
			if lerr == false {
				context.String(404, "404 page not found")
				context.Abort()
			} else {
				context.HTMLString(http.StatusInternalServerError, `<html>
				<head><title>Syntax Error in `+context.Request.URL.Path+`</title>
				<body>
					<h1>Syntax Error in file `+context.Request.URL.Path+`</h1>
					<code>`+string(err.Error())+`</code>
				</body>
				</html>`)
				context.Abort()
			}
		}
		L.LoadBuffer(code, len(code), file) // This shouldn't error, was checked earlier.
		if L.Pcall(0, 0, 0) != 0 {          // != 0 means error in execution
			context.HTMLString(http.StatusInternalServerError, `<html>
			<head><title>Runtime Error in `+context.Request.URL.Path+`</title>
			<body>
				<h1>Runtime Error in file `+context.Request.URL.Path+`</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
		}
		/*L.DoString("return CONTENT_TO_RETURN")
		v := luar.CopyTableToMap(L, nil, -1)
		m := v.(map[string]interface{})
		i := int(m["code"].(float64))
		if err != nil {
			i = http.StatusOK
		}*/
		//context.HTMLString(i, m["content"].(string))
	}
}

func main() {
	cbc = cache.New(5*time.Minute, 30*time.Second) // Initialize cache with 5 minute lifetime and purge every 30 seconds
	cfe = cache.New(5*time.Minute, 30*time.Second) // File-Exists Cache

	// Use config
	flag.String("config", "", "Parse Config File")

	var host = flag.String("host", "", "IP of Host to bind the Webserver on")
	var port = flag.Int("port", 8080, "Port to run Webserver on")
	jobs = flag.Int("states", 16, "Number of Preinitialized Lua States")
	var workers = flag.Int("workers", runtime.NumCPU(), "Number of Worker threads.")
	var webroot = flag.String("root", ".", "Path to Web Root")

	// Middleware options
	useRecovery := flag.Bool("recovery", false, "Recover from Panics")
	useLogger := flag.Bool("logger", true, "Log Requests and Cache information")
	useGzip := flag.Bool("gzip", false, "Use GZIP")

	flag.Parse()

	runtime.GOMAXPROCS(*workers)

	srv := new_server()
	if *useLogger {
		doLog = true
		srv.Use(gin.Logger())
	}
	if *useRecovery {
		srv.Use(gin.Recovery())
	}
	if *useGzip {
		srv.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	bootstrap(srv, *webroot)

	srv.Run(*host + ":" + strconv.Itoa(*port))
}
