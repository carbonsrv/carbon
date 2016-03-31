package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/carbonsrv/carbon/modules/helpers"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pmylund/go-cache"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
)

// Cache
var kvstore *cache.Cache
var cbc *cache.Cache
var cfe *cache.Cache
var LDumper *lua.State

// cacheDump dumps the source bytecode cached
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

// bcdump actually dumps the bytecode
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

// Preloaded is the chan that contains the preloaded states
var Preloaded chan *lua.State

// Preloader is the function that preloads the states
func Preloader() {
	Preloaded = make(chan *lua.State, jobs)
	for {
		//fmt.Println("preloading")
		L := luar.Init()
		Bind(L, webroot)
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

// GetInstance grabs an instance from the preloaded list
func GetInstance() *lua.State {
	//fmt.Println("grabbing instance")
	L := <-Preloaded
	//fmt.Println("Done")
	return L
}

// Init
func Init(j int, cfe_new *cache.Cache, kvstore_new *cache.Cache, root string) {
	webroot = root
	cfe = cfe_new
	kvstore = kvstore_new
	jobs = j
	filesystem = physfs.FileSystem()
	cbc = cache.New(5*time.Minute, 30*time.Second) // Initialize cache with 5 minute lifetime and purge every 30 seconds
	LDumper = luar.Init()
}

// Lua behaves PHP-like, but for Lua
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
	<pre>`+string(err.Error())+`</pre>
</body>
</html>`)
					context.Abort()
					L.Close()
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
	<pre>`+L.ToString(-1)+`</pre>
</body>
</html>`)
				context.Abort()
				L.Close()
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
			L.Close()
		} else {
			context.Next()
		}
	}
}

// DLR_NS does Route creation by lua
func DLR_NS(bcode string, dobind bool, vals map[string]interface{}) (func(*gin.Context), error) {
	/*code, err := bcdump(code)
	if err != nil {
		return func(*gin.Context) {}, err
	}*/
	return func(context *gin.Context) {
		L := GetInstance()
		if dobind {
			luar.Register(L, "", vals)
		}
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
	<pre>`+L.ToString(-1)+`</pre>
</body>
</html>`)
			context.Abort()
			L.Close()
			return
		}
		L.Close()
	}, nil
}

// DLR_RUS does the same as above, but reuses states
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
			luar.Register(L, "", vals)
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
	<pre>`+L.ToString(-1)+`</pre>
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
	CheckOrigin: func(r *http.Request) bool { // Because it breaks some things.
		return true
	},
}

// DLRWS_NS does the same as the first one, but for websockets.
func DLRWS_NS(bcode string, dobind bool, vals map[string]interface{}) (func(*gin.Context), error) { // Same as above, but for websockets.
	return func(context *gin.Context) {
		L := GetInstance()
		BindContext(L, context)
		if dobind {
			luar.Register(L, "", vals)
		}

		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			fmt.Println("Websocket error: " + err.Error()) // silent error.
		}
		luar.Register(L, "ws", luar.Map{
			"con":           conn,
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
			"read_con": (func(con *websocket.Conn) (int, string, error) {
				messageType, p, err := con.ReadMessage()
				if err != nil {
					return -1, "", err
				}
				return messageType, string(p), nil
			}),
			"send": (func(t int, cnt string) error {
				return conn.WriteMessage(t, []byte(cnt))
			}),
			"send_con": (func(con *websocket.Conn, t int, cnt string) error {
				return con.WriteMessage(t, []byte(cnt))
			}),
			"close": (func() error {
				return conn.Close()
			}),
			"close_con": (func(con *websocket.Conn) error {
				return con.Close()
			}),
		})
		L.LoadBuffer(bcode, len(bcode), "route")
		if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
			fmt.Println("Websocket Lua error: " + L.ToString(-1))
			context.Abort()
			return
		}
		// Close websocket.
		conn.Close()

		L.Close()
	}, nil
}

// DLRWS_RUS also does the thing for websockets, but reuses states, not quite handy given how many connections a websocket could take and how long the connection could keep alive.
func DLRWS_RUS(bcode string, instances int, dobind bool, vals map[string]interface{}) (func(*gin.Context), error) { // Same as above, but reusing states.
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
			luar.Register(L, "", vals)
		}
		schan <- L
	}
	return func(context *gin.Context) {
		L := <-schan
		BindContext(L, context)
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			fmt.Println("Websocket error: " + err.Error()) // silent error.
		}
		luar.Register(L, "ws", luar.Map{
			"con":           conn,
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
			"read_con": (func(con *websocket.Conn) (int, string, error) {
				messageType, p, err := con.ReadMessage()
				if err != nil {
					return -1, "", err
				}
				return messageType, string(p), nil
			}),
			"send": (func(t int, cnt string) error {
				return conn.WriteMessage(t, []byte(cnt))
			}),
			"send_con": (func(con *websocket.Conn, t int, cnt string) error {
				return con.WriteMessage(t, []byte(cnt))
			}),
			"close": (func() error {
				return conn.Close()
			}),
			"close_con": (func(con *websocket.Conn) error {
				return con.Close()
			}),
		})
		L.LoadBuffer(bcode, len(bcode), "route")
		if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
			fmt.Println("Websocket Lua error: " + L.ToString(-1))
			context.Abort()
			return
		}
		schan <- L
		// Close websocket.
		conn.Close()
	}, nil
}
