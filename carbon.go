package main

//go:generate go-bindata -nomemcopy -o modules/glue/generated_lua.go -pkg=glue -prefix "./builtin" ./builtin ./builtin/3rdparty ./builtin/libs ./builtin/libs/wrappers

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/carbonsrv/carbon/modules/luaconf"
	"github.com/carbonsrv/carbon/modules/middleware"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
	"github.com/pmylund/go-cache"
	"github.com/vifino/golua/lua"
	"golang.org/x/net/http2"
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

// File system functions
var filesystem http.FileSystem

func initPhysFS(path string) http.FileSystem {
	err := physfs.Init()
	if err != nil {
		panic(err)
	}
	err = physfs.Mount(path, "/root", true)
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
	file = "/root/" + file
	if physfs.Exists(file) {
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
	return "", errors.New("no such file or directory.")
}
func fileExists(file string) bool {
	return physfs.Exists("/root/" + file)
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
var timeout time.Duration = 10 * time.Second

func new_server() *gin.Engine {
	return gin.New()
}

func serve(srv http.Handler, en_http bool, en_https bool, en_http2 bool, bind string, binds string, cert string, key string) {
	end := make(chan bool)
	if en_http {
		go serveHTTP(srv, bind, en_http2)
	}
	if en_https {
		cert, _ := filepath.Abs(cert)
		key, _ := filepath.Abs(key)
		go serveHTTPS(srv, binds, en_http2, cert, key)
	}
	<-end
}
func serveHTTP(srv http.Handler, bind string, en_http2 bool) {
	s := &http.Server{
		Addr:           bind,
		Handler:        srv,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
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
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
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
	var eval = flag.String("eval", "", "Eval Lua Code")
	var is_app = flag.Bool("app", false, "Script is bundle")
	var licenses = flag.Bool("licenses", false, "Show licenses")

	var host = flag.String("host", "", "IP of Host to bind the Webserver on")
	var port = flag.Int("port", 80, "Port to run Webserver on (HTTP)")
	var ports = flag.Int("ports", 443, "Port to run Webserver on (HTTPS)")
	var cert = flag.String("cert", "", "Certificate File for HTTPS")
	var key = flag.String("key", "", "Key File for HTTPS")
	var en_http = flag.Bool("http", true, "Listen HTTP")
	var en_https = flag.Bool("https", false, "Listen HTTPS")
	var en_http2 = flag.Bool("http2", false, "Enable HTTP/2")
	var set_timeout = flag.Int64("timeout", 0, "Timeout for HTTP read/write calls. (Seconds)")

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

	flag.Parse()

	timeout = time.Duration(*set_timeout) * time.Second
	if *en_https {
		if *key == "" || *cert == "" {
			panic("Need to have a Key and a Cert defined.")
		}
	}

	if *licenses {
		fmt.Println(glue.GetGlue("NOTICE.txt"))
		os.Exit(0)
	}

	var script string
	if *script_flag == "" {
		script = flag.Arg(0)
	} else {
		script = *script_flag
	}

	runtime.GOMAXPROCS(*workers)

	root, _ := filepath.Abs(*webroot)
	physroot_path := root
	if *is_app {
		physroot_path = script
	}
	filesystem = initPhysFS(physroot_path)

	defer physfs.Deinit()
	go scheduler.Run()                         // Run the scheduler.
	go middleware.Preloader()                  // Run the Preloader.
	middleware.Init(*jobs, cfe, kvstore, root) // Run init sequence.

	if *doDebug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	args := flag.Args()

	bindhook := func(_ *lua.State) {}

	if *eval != "" {
		err := luaconf.Eval(*eval, args, cfe, root, *useRecovery, *useLogger, *run_repl, func(srv *gin.Engine) {
			serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
		}, bindhook)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if script == "" {
		if *run_repl {
			err := luaconf.REPL(args, cfe, root, *useRecovery, false, func(srv *gin.Engine) {
				serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
			}, bindhook)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err := luaconf.Eval(glue.GetGlue("bootscript.lua"), args, cfe, root, *useRecovery, *useLogger, *run_repl, func(srv *gin.Engine) {
				serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
			}, bindhook)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	} else {
		if *useLogger {
			doLog = true
		}

		if flag.Arg(1) == "" {
			args = make([]string, 0)
		} else {
			args = args[1:]
		}
		if *is_app {
			script_src, err := fileRead("app.lua")
			if err != nil {
				script_base := filepath.Base(script)
				script_name := script_base[:len(script_base)-len(filepath.Ext(script))]
				script_src, err = fileRead(script_name + ".lua")
				if err != nil {
					script_src, err = fileRead("init.lua")
					if err != nil {
						fmt.Println("Error: App Bundle does not contain 'app.lua', '" + script_name + ".lua' or 'init.lua'. No idea what to run. Aborting.")
						os.Exit(1)
					}
				}
			}
			err = luaconf.Eval(script_src, args, cfe, script, *useRecovery, *useLogger, *run_repl, func(srv *gin.Engine) {
				serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
			}, bindhook)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err := luaconf.Configure(script, args, cfe, root, *useRecovery, *useLogger, *run_repl, func(srv *gin.Engine) {
				serve(srv, *en_http, *en_https, *en_http2, *host+":"+strconv.Itoa(*port), *host+":"+strconv.Itoa(*ports), *cert, *key)
			}, bindhook)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
