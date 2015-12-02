package main

//go:generate go-bindata -o modules/glue/generated_lua.go -pkg=glue -prefix "./lua" ./lua ./lua/3rdparty ./lua/libs ./lua/libs/wrappers

import (
	"bufio"
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/bradfitz/http2"
	"github.com/carbonsrv/carbon/modules/luaconf"
	"github.com/carbonsrv/carbon/modules/middleware"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/carbonsrv/carbon/modules/static"
	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
	"github.com/pmylund/go-cache"
	"github.com/vifino/contrib/gzip"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// General
var jobs *int

// Cache
var cfe *cache.Cache
var kvstore *cache.Cache

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
		"***":  staticServe.ServeCached("", staticServe.PhysFS("", "", true, true), cfe),
	})
	/*srv.GET(`/:file`, switcher)
	srv.POST(`/:file`, switcher)*/
	srv.Use(switcher)
	/*st := staticServe.ServeCached("", staticServe.PhysFS("", true, true), cfe)
	lr := luaroute(dir)
	srv.Use(lr)
	srv.Use(st)*/
}

func serve(srv http.Handler, en_http bool, en_https bool, en_http2 bool, bind string, binds string, cert string, key string) {
	end := make(chan bool)
	if en_http {
		go serveHTTP(srv, bind, en_http2)
	}
	if en_https {
		go serveHTTPS(srv, binds, en_http2, cert, key)
	}
	<-end
}
func serveHTTP(srv http.Handler, bind string, en_http2 bool) {
	s := &http.Server{
		Addr:           bind,
		Handler:        srv,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if en_http2 {
		http2.ConfigureServer(s, nil)
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
func serveHTTPS(srv http.Handler, bind string, en_http2 bool, cert string, key string) {
	s := &http.Server{
		Addr:           bind,
		Handler:        srv,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if en_http2 {
		http2.ConfigureServer(s, nil)
	}
	err := s.ListenAndServeTLS(cert, key)
	if err != nil {
		panic(err)
	}
}

func main() {
	cfe = cache.New(5*time.Minute, 30*time.Second) // File-Exists Cache
	kvstore = cache.New(-1, -1)                    // Key-Value Storage

	// Use config
	flag.String("config", "", "Parse Config File")
	var script_flag = flag.String("script", "", "Parse Lua Script as initialization")
	var run_repl = flag.Bool("repl", false, "Run REPL")

	var host = flag.String("host", "", "IP of Host to bind the Webserver on")
	var port = flag.Int("port", 8080, "Port to run Webserver on (HTTP)")
	var ports = flag.Int("ports", 8443, "Port to run Webserver on (HTTPS)")
	var cert = flag.String("cert", "", "Certificate File for HTTPS")
	var key = flag.String("key", "", "Key File for HTTPS")
	var en_http = flag.Bool("http", true, "Listen HTTP")
	var en_https = flag.Bool("https", false, "Listen HTTPS")
	var en_http2 = flag.Bool("http2", false, "Enable HTTP/2")

	wrkrs := 2
	if runtime.NumCPU() > 2 {
		wrkrs = runtime.NumCPU()
	}
	jobs = flag.Int("states", wrkrs, "Number of Preinitialized Lua States")
	var workers = flag.Int("workers", wrkrs, "Number of Worker threads.")
	var webroot = flag.String("root", ".", "Path to Web Root")

	// Do debug!
	doDebug := flag.Bool("debug", false, "Show debug information")
	// Middleware options
	useRecovery := flag.Bool("recovery", false, "Recover from Panics")
	useLogger := flag.Bool("logger", true, "Log Request information")
	useGzip := flag.Bool("gzip", false, "Use GZIP")

	flag.Parse()

	if *en_https {
		if *key == "" || *cert == "" {
			panic("Need to have a Key and a Cert defined.")
		}
	}

	runtime.GOMAXPROCS(*workers)

	root, _ := filepath.Abs(*webroot)
	filesystem = initPhysFS(root)
	defer physfs.Deinit()
	go scheduler.Run()                   // Run the scheduler.
	go middleware.Preloader()            // Run the Preloader.
	middleware.Init(*jobs, cfe, kvstore) // Run init sequence.

	var script string
	if *script_flag == "" {
		script = flag.Arg(0)
	} else {
		script = *script_flag
	}

	if *doDebug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	if script == "" {
		if *run_repl {
			err := luaconf.REPL(cfe, root, *useRecovery, *useLogger, func(srv *gin.Engine) {
				serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
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
			serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
		}
	} else {
		if *useLogger {
			doLog = true
		}
		err := luaconf.Configure(script, cfe, root, *useRecovery, *useLogger, func(srv *gin.Engine) {
			serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
