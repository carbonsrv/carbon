package routes

import (
	"../glue"
	"../scheduler"
	"bufio"
	"errors"
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"net/http"
	"time"
)

// Cache
var cbc *cache.Cache
var LDumper *lua.State

func cacheDump(file string) (string, error, bool) {
	data_tmp, found := cbc.Get(file)
	if found == false {
		data, err := fileRead(file)
		if err != nil {
			return "", err, false
		}
		res, err := bcdump(data)
		if err != nil {
			return "", err, true
		}
		cbc.Set(file, res, cache.DefaultExpiration)
		return res, nil, false
	} else {
		//debug("Using Bytecode-cache for " + file)
		return data_tmp.(string), nil, false
	}
}
func bcdump(data string) (string, error) {
	if LDumper.LoadString(data) != 0 {
		return "", errors.New(LDumper.ToString(-1))
	}
	defer LDumper.Pop(1)
	return LDumper.FDump(), nil
}

// FS
var filesystem http.FileSystem

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

// Preloader/Starter
var jobs int
var Preloaded chan *lua.State

func Preloader() {
	Preloaded = make(chan *lua.State, jobs)
	for {
		//fmt.Println("preloading")
		state := luar.Init()
		err := state.DoString(glue.MainGlue())
		if err != nil {
			fmt.Println(err)
		}
		Preloaded <- state
	}
}
func GetInstance() *lua.State {
	//fmt.Println("grabbing instance")
	L := <-Preloaded
	//fmt.Println("Done")
	return L
}

// Init
func Init(j int) {
	jobs = j
	filesystem = physfs.FileSystem()
	cbc = cache.New(5*time.Minute, 30*time.Second) // Initialize cache with 5 minute lifetime and purge every 30 seconds
	LDumper = luar.Init()
}

// PHP-like lua scripts
func Lua() func(*gin.Context) {
	//LDumper := luar.Init()
	return func(context *gin.Context) {
		//fmt.Println("start")
		L := GetInstance()
		//fmt.Println("after start")
		defer scheduler.Add(func() {
			L.Close()
		})
		file := context.Request.URL.Path
		//fmt.Println("after after start")
		luar.Register(L, "", luar.Map{
			"context": context,
			"req":     context.Request,
		})
		//fmt.Println("before cache")
		code, err, lerr := cacheDump(file)
		//fmt.Println("after cache")
		if err != nil {
			if lerr == false {
				context.Next()
				return
			} else {
				context.HTMLString(http.StatusInternalServerError, `<html>
				<head><title>Syntax Error in `+context.Request.URL.Path+`</title>
				<body>
					<h1>Syntax Error in file `+context.Request.URL.Path+`:</h1>
					<code>`+string(err.Error())+`</code>
				</body>
				</html>`)
				context.Abort()
				return
			}
		}
		//fmt.Println("before loadbuffer")
		L.LoadBuffer(code, len(code), file) // This shouldn't error, was checked earlier.
		if L.Pcall(0, 0, 0) != 0 {          // != 0 means error in execution
			context.HTMLString(http.StatusInternalServerError, `<html>
			<head><title>Runtime Error in `+context.Request.URL.Path+`</title>
			<body>
				<h1>Runtime Error in file `+context.Request.URL.Path+`:</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
			return
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

// Route creation by lua
func New(bcode string) (func(*gin.Context), error) {
	/*code, err := bcdump(code)
	if err != nil {
		return func(*gin.Context) {}, err
	}*/
	return func(context *gin.Context) {
		L := GetInstance()
		L.DoString(glue.RouteGlue())
		luar.Register(L, "", luar.Map{
			"context": context,
			"req":     context.Request,
		})
		/*luar.Register(L, "fs", luar.Map{ // PhysFS
			"mount":       physfs.Mount,
			"exits":       physfs.Exists,
			"getFS":       physfs.FileSystem,
			"mkdir":       physfs.Mkdir,
			"umount":      physfs.RemoveFromSearchPath,
			"delete":      physfs.Delete,
			"setWriteDir": physfs.SetWriteDir,
			"getWriteDir": physfs.GetWriteDir,
		})
		luar.Register(L, "mw", luar.Map{
			"Lua": Lua,
			"ExtRoute": (func(plan map[string]interface{}) func(*gin.Context) {
				newplan := make(Plan, len(plan))
				for k, v := range plan {
					newplan[k] = v.(func(*gin.Context))
				}
				return ExtRoute(newplan)
			}),
			"Logger":   gin.Logger,
			"Recovery": gin.Recovery,
			"GZip": func() func(*gin.Context) {
				return gzip.Gzip(gzip.DefaultCompression)
			},
			"New": New,
		})*/
		//fmt.Println("before loadbuffer")
		if L.LoadBuffer(bcode, len(bcode), "route") != 0 {
			context.HTMLString(http.StatusInternalServerError, `<html>
			<head><title>Syntax Error in `+context.Request.URL.Path+`</title>
			<body>
				<h1>Syntax Error in Lua Route on `+context.Request.URL.Path+`:</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
			return
		}
		if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
			context.HTMLString(http.StatusInternalServerError, `<html>
			<head><title>Runtime Error on `+context.Request.URL.Path+`</title>
			<body>
				<h1>Runtime Error in Lua Route on `+context.Request.URL.Path+`:</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
			return
		}
	}, nil
}
