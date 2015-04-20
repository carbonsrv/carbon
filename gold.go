package main

//go:generate go-bindata -o glue/generated_lua.go -pkg=glue -prefix "./lua" ./lua

import (
	"./glue"
	"fmt"
	//"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"io/ioutil"
	//"net/http"
	//"runtime"
	"log"
	"os"
	"strings"
	"time"
)

// Preloader/Starter
var jobs int = 32
var preloaded chan *lua.State

func preloader() {
	preloaded = make(chan *lua.State, jobs)
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
		log.Printf("Using cache for %s", file)
		res = data_tmp.(string)
	}
	return res, nil
}
func cacheFileExists(c *cache.Cache, file string) bool {
	data_tmp, found := c.Get(file)
	if found == false {
		exists := fileExists(file)
		c.Set(file, exists, cache.DefaultExpiration)
		return exists
	} else {
		return data_tmp.(bool)
	}
}

// FS Helper
func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func new_server() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	return r
}

// Routes
func logic_switcher(dir string) func(*gin.Context) {
	c := cache.New(5*time.Minute, 30*time.Second) // File-Exists Cache
	return func(context *gin.Context) {
		file := context.Params.ByName("file")
		fe := cacheFileExists(c, file)
		if fe == true {
			if strings.HasSuffix(file, ".lua") {
				luaroute(dir)(context)
			} else {
				static.Serve("/", static.LocalFile(dir, false))(context)
			}
		} else {
			context.String(404, "404 page not found")
		}
	}
}

func luaroute(dir string) func(*gin.Context) {
	c := cache.New(5*time.Minute, 30*time.Second) // Initialize cache with 5 minute lifetime and purge every 30 seconds
	return func(context *gin.Context) {
		L := getInstance()
		file := dir + context.Request.URL.Path
		luar.Register(L, "", luar.Map{
			"context": context,
			"req":     context.Request,
		})
		code, err := cacheRead(c, file)
		if err != nil {
			context.String(404, "404 page not found")
		}
		err = L.DoString(code)
		if err != nil {
			context.HTMLString(500, `<html>
			<head><title>Error in `+context.Request.URL.Path+`</title>
			<body>
				<h1>Error in file `+context.Request.URL.Path+`</h1>
				<code>`+string(err.Error())+`</code>
			</body>
			</html>`)
		}
		L.DoString("return CONTENT_TO_RETURN")
		v := luar.CopyTableToMap(L, nil, -1)
		m := v.(map[string]interface{})
		i := int(m["code"].(float64))
		if err != nil {
			i = 200
		}
		defer L.Close()
		context.HTMLString(i, m["content"].(string))
	}
}

func run(dir string) {
	go preloader() // Run the instance starter.
	srv := new_server()
	//srv.Use(gzip.Gzip(gzip.DefaultCompression))
	srv.GET(`/:file`, logic_switcher(dir))

	//srv.Use(martini.Static(dir))
	srv.Run(":3000")
}

func main() {
	run(".")
}
