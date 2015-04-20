package main

//go:generate go-bindata -o glue/generated_lua.go -pkg=glue -prefix "./lua" ./lua

import (
	"./glue"
	"fmt"
	martini "github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"io/ioutil"
	"net/http"
	"strconv"
)

var jobs int = 8
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
	go preloader()
	srv := new_server(dir)
	srv.Use(gzip.All())
	srv.Get("/**.lua", func(res http.ResponseWriter, req *http.Request) (int, string) {
		L := getInstance()
		res.Header().Set("Content-Type", "text/html")
		file := dir + req.URL.Path
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(file)
			return 404, "404 page not found"
		}
		code := string(data)
		luar.Register(L, "", luar.Map{
			"res": res,
			"req": req,
		})
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
	srv.Get("/**", martini.Static(dir))
	//srv.Use(martini.Static(dir))
	srv.Run()
}

func main() {
	run(".")
}
