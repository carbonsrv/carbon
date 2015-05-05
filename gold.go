package main

//go:generate go-bindata -o glue/generated_lua.go -pkg=glue -prefix "./lua" ./lua

import (
	"./glue"
	"./modules/static"
	"./scheduler"
	"bufio"
	"errors"
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
	"github.com/pmylund/go-cache"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"log"
	"net/http"
	"path/filepath"
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
		data, err := fileRead(file)
		if err != nil {
			return "", err
		}
		c.Set(file, data, cache.DefaultExpiration)
	} else {
		debug("Using cache for %s" + file)
		res = data_tmp.(string)
	}
	return res, nil
}
func cacheDump(L *lua.State, file string) (string, error, bool) {
	data_tmp, found := cbc.Get(file)
	if found == false {
		data, err := fileRead(file)
		if err != nil {
			return "", err, false
		}
		if L.LoadString(data) != 0 {
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

// File sytem functions
var filesystem http.FileSystem

func initPhysFS(path string) http.FileSystem {
	err := physfs.Init()
	if err != nil {
		panic(err)
	}
	err = physfs.Mount(path, "/", true)
	if err != nil {
		panic(err)
	}
	return physfs.FileSystem()
}

/*func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}*/
func fileRead(file string) (string, error) {
	f, err := filesystem.Open(file)
	defer f.Close()
	if err != nil {
		return "", err
	}
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	r := bufio.NewReader(f)
	buf := make([]byte, fi.Size())
	_, err = r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), err
}
func fileExists(file string) bool {
	return physfs.Exists(file)
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

// Logging
var doLog bool = false

func debug(str string) {
	if doLog {
		log.Print(str)
	}
}

// Server
func new_server() *gin.Engine {
	return gin.New()
}
func bootstrap(srv *gin.Engine, dir string) *gin.Engine {
	go preloader()     // Run the instance starter.
	go scheduler.Run() // Run the scheduler.
	switcher := logic_switcheroo(dir)
	/*srv.GET(`/:file`, switcher)
	srv.POST(`/:file`, switcher)*/
	srv.Use(switcher)
	return srv
}

// Routes
func logic_switcheroo(dir string) func(*gin.Context) {
	st := staticServe.ServeCached("", staticServe.PhysFS("", true, true))
	lr := luaroute(dir)
	return func(context *gin.Context) {
		file := dir + context.Request.URL.Path
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
	root, _ := filepath.Abs(*webroot)
	debug(root)
	filesystem = initPhysFS(root)
	defer physfs.Deinit()
	bootstrap(srv, "")

	srv.Run(*host + ":" + strconv.Itoa(*port))
}
