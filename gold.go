package main

//go:generate go-bindata -o modules/glue/generated_lua.go -pkg=glue -prefix "./lua" ./lua

import (
	"./modules/luaconf"
	"./modules/middleware"
	"./modules/scheduler"
	"./modules/static"
	"bufio"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
	"github.com/pmylund/go-cache"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// General
var jobs *int

// Cache
var cfe *cache.Cache

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
func bootstrap(srv *gin.Engine, dir string, cfe *cache.Cache) {
	switcher := middleware.ExtRoute(middleware.Plan{
		".lua": middleware.Lua(),
		"***":  staticServe.ServeCached("", staticServe.PhysFS("", true, true), cfe),
	})
	/*srv.GET(`/:file`, switcher)
	srv.POST(`/:file`, switcher)*/
	srv.Use(switcher)
	/*st := staticServe.ServeCached("", staticServe.PhysFS("", true, true), cfe)
	lr := luaroute(dir)
	srv.Use(lr)
	srv.Use(st)*/
}

func main() {
	cfe = cache.New(5*time.Minute, 30*time.Second) // File-Exists Cache

	// Use config
	flag.String("config", "", "Parse Config File")
	var script = flag.String("script", "", "Parse Lua Script as initialization")

	var host = flag.String("host", "", "IP of Host to bind the Webserver on")
	var port = flag.Int("port", 8080, "Port to run Webserver on (HTTP)")
	jobs = flag.Int("states", runtime.NumCPU(), "Number of Preinitialized Lua States")
	var workers = flag.Int("workers", runtime.NumCPU(), "Number of Worker threads.")
	var webroot = flag.String("root", ".", "Path to Web Root")

	// Middleware options
	useRecovery := flag.Bool("recovery", false, "Recover from Panics")
	useLogger := flag.Bool("logger", true, "Log Requests and Cache information")
	useGzip := flag.Bool("gzip", false, "Use GZIP")

	flag.Parse()

	runtime.GOMAXPROCS(*workers)

	root, _ := filepath.Abs(*webroot)
	filesystem = initPhysFS(root)
	defer physfs.Deinit()
	go scheduler.Run()        // Run the scheduler.
	go middleware.Preloader() // Run the Preloader.
	middleware.Init(*jobs)    // Run init sequence.

	if *script == "" {
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
		bootstrap(srv, "", cfe)
		srv.Run(*host + ":" + strconv.Itoa(*port))
	} else {
		srv, err := luaconf.Configure(*script, cfe, *webroot)
		if err != nil {
			panic(err)
		}
		srv.Run(*host + ":" + strconv.Itoa(*port))
	}
}
