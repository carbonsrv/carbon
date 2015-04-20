package main

//go:generate go-bindata -o glue/generated_lua.go -pkg=glue -prefix "./lua" ./lua

import (
	"./glue"
	"fmt"
	martini "github.com/go-martini/martini"
	//"github.com/martini-contrib/gzip"
	"github.com/pmylund/go-cache"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"io/ioutil"
	"net/http"
	//"runtime"
	"log"
	"strconv"
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

func new_server(dir string) *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	//m.Use(martini.Static(dir))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}

func run(dir string) {
	c := cache.New(5*time.Minute, 30*time.Second) // Initialize cache with 5 minute lifetime and purge every 30 seconds
	go preloader()                                // Run the instance starter.
	srv := new_server(dir)
	//srv.Use(gzip.All())
	srv.Get("/**.lua", func(res http.ResponseWriter, req *http.Request) (int, string) {
		L := getInstance()
		res.Header().Set("Content-Type", "text/html")
		file := dir + req.URL.Path
		luar.Register(L, "", luar.Map{
			"res": res,
			"req": req,
		})
		code, err := cacheRead(c, file)
		if err != nil {
			return 404, "404 page not found"
		}
		err = L.DoString(code)
		if err != nil {
			return 500, `<html>
			<head><title>Error in "` + req.URL.Path + `"</title>
			<body>
				<h1>Error in file "` + req.URL.Path + `"</h1>
				<code>` + string(err.Error()) + `</code>
			</body>
			</html>`
		}
		L.DoString("return CONTENT_TO_RETURN")
		v := luar.CopyTableToMap(L, nil, -1)
		m := v.(map[string]interface{})
		i, err := strconv.Atoi(m["code"].(string))
		if err != nil {
			i = 200
		}
		defer L.Close()
		return i, m["content"].(string)
	})
	srv.Get("/:file", martini.Static(dir))
	//srv.Use(martini.Static(dir))
	srv.Run()
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	martini.Env = martini.Prod
	run(".")
}
