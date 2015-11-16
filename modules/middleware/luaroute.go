package middleware

import (
	"bufio"
	"errors"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pmylund/go-cache"
	"github.com/vifino/carbon/modules/glue"
	"github.com/vifino/carbon/modules/helpers"
	"github.com/vifino/carbon/modules/scheduler"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"net/http"
	"time"
)

// Cache
var cbc *cache.Cache
var cfe *cache.Cache
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

func fileExists(file string) bool {
	data_tmp, found := cfe.Get(file)
	if found == false {
		exists := physfs.Exists(file)
		cfe.Set(file, exists, cache.DefaultExpiration)
		return exists
	} else {
		return data_tmp.(bool)
	}
}

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
		L := luar.Init()
		Bind(L)
		err := L.DoString(glue.MainGlue())
		if err != nil {
			panic(err)
		}
		err = L.DoString(glue.RouteGlue())
		if err != nil {
			panic(err)
		}
		Preloaded <- L
	}
}
func GetInstance() *lua.State {
	//fmt.Println("grabbing instance")
	L := <-Preloaded
	//fmt.Println("Done")
	return L
}

// Init
func Init(j int, cfe_new *cache.Cache) {
	cfe = cfe_new
	jobs = j
	filesystem = physfs.FileSystem()
	cbc = cache.New(5*time.Minute, 30*time.Second) // Initialize cache with 5 minute lifetime and purge every 30 seconds
	LDumper = luar.Init()
}

// PHP-like lua scripts
func Lua() func(*gin.Context) {
	//LDumper := luar.Init()
	return func(context *gin.Context) {
		file := context.Request.URL.Path
		if fileExists(file) {
			//fmt.Println("start")
			L := GetInstance()
			//fmt.Println("after start")
			defer scheduler.Add(func() {
				L.Close()
			})
			//fmt.Println("after after start")
			BindContext(L, context)
			//fmt.Println("before cache")
			code, err, lerr := cacheDump(file)
			//fmt.Println("after cache")
			if err != nil {
				if lerr == false {
					context.Next()
					return
				} else {
					helpers.HTMLString(context, http.StatusInternalServerError, `<html>
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
				helpers.HTMLString(context, http.StatusInternalServerError, `<html>
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
			//helpers.HTMLString(context, i, m["content"].(string))
		} else {
			context.Next()
		}
	}
}

// Route creation by lua
func DLR_NS(bcode string, dobind bool, vals map[string]interface{}) (func(*gin.Context), error) {
	/*code, err := bcdump(code)
	if err != nil {
		return func(*gin.Context) {}, err
	}*/
	return func(context *gin.Context) {
		L := GetInstance()
		if dobind {
			for k, v := range vals {
				luar.Register(L, "", luar.Map{
					k: v,
				})
			}
		}
		defer scheduler.Add(func() {
			L.Close()
		})
		BindContext(L, context)
		//fmt.Println("before loadbuffer")
		/*if L.LoadBuffer(bcode, len(bcode), "route") != 0 {
			helpers.HTMLString(context, http.StatusInternalServerError, `<html>
			<head><title>Syntax Error in `+context.Request.URL.Path+`</title>
			<body>
				<h1>Syntax Error in Lua Route on `+context.Request.URL.Path+`:</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
			return
		}*/
		L.LoadBuffer(bcode, len(bcode), "route")
		if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
			helpers.HTMLString(context, http.StatusInternalServerError, `<html>
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
func DLR_RUS(bcode string, instances int, dobind bool, vals map[string]interface{}) (func(*gin.Context), error) { // Same as above, but reuses states. Much faster. Higher memory use though, because more states.
	insts := instances
	if instances < 0 {
		insts = 2
		if jobs/2 > 1 {
			insts = jobs
		}
	}
	schan := make(chan *lua.State, insts)
	for i := 0; i < jobs/2; i++ {
		L := GetInstance()
		if dobind {
			for k, v := range vals {
				luar.Register(L, "", luar.Map{
					k: v,
				})
			}
		}
		if L.LoadBuffer(bcode, len(bcode), "route") != 0 {
			return func(context *gin.Context) {}, errors.New(L.ToString(-1))
		}
		L.PushValue(-1)
		schan <- L
	}
	return func(context *gin.Context) {
		L := <-schan
		BindContext(L, context)
		if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
			helpers.HTMLString(context, http.StatusInternalServerError, `<html>
			<head><title>Runtime Error on `+context.Request.URL.Path+`</title>
			<body>
				<h1>Runtime Error in Lua Route on `+context.Request.URL.Path+`:</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
			return
		}
		L.PushValue(-1)
		schan <- L
	}, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func DLRWS_RUS(bcode string, instances int, dobind bool, vals map[string]interface{}) (func(*gin.Context), error) { // Same as above, but for websockets. Not working because?!
	insts := instances
	if instances < 0 {
		insts = 2
		if jobs/2 > 1 {
			insts = jobs
		}
	}
	schan := make(chan *lua.State, insts)
	for i := 0; i < jobs/2; i++ {
		L := GetInstance()
		if dobind {
			for k, v := range vals {
				luar.Register(L, "", luar.Map{
					k: v,
				})
			}
		}
		if L.LoadBuffer(bcode, len(bcode), "route") != 0 {
			return func(context *gin.Context) {}, errors.New(L.ToString(-1))
		}
		L.PushValue(-1)
		schan <- L
	}
	return func(context *gin.Context) {
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			return // silent error.
		}
		L := <-schan
		BindContext(L, context)
		luar.Register(L, "ws", luar.Map{
			"BinaryMessage": websocket.BinaryMessage,
			"TextMessage":   websocket.TextMessage,
			//"read":          conn.ReadMessage,
			//"send":          conn.SendMessage,
			"read": (func() (int, string, error) {
				messageType, p, err := conn.ReadMessage()
				if err != nil {
					return -1, "", err
				}
				return messageType, string(p), nil
			}),
			"send": (func(t int, cnt string) error {
				return conn.WriteMessage(t, []byte(cnt))
			}),
		})
		if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
			helpers.HTMLString(context, http.StatusInternalServerError, `<html>
			<head><title>Runtime Error on `+context.Request.URL.Path+`</title>
			<body>
				<h1>Runtime Error in Lua Route on `+context.Request.URL.Path+`:</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
			return
		}
		L.PushValue(-1)
		schan <- L
	}, nil
}
