// Big monolithic binding file.
// Binds a ton of things.

package middleware

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/GeertJohan/go.linenoise"
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/carbonsrv/carbon/modules/helpers"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/carbonsrv/carbon/modules/static"
	"github.com/fzzy/radix/redis"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/vifino/contrib/gzip"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// Vars
var webroot string

// Bind all things
func Bind(L *lua.State, root string) {
	webroot = root

	luar.Register(L, "var", luar.Map{ // Vars
		"root": root,
	})

	BindCarbon(L)
	BindMiddleware(L)
	BindRedis(L)
	BindKVStore(L)
	BindPhysFS(L)
	BindIOEnhancements(L)
	BindOSEnhancements(L)
	BindThread(L)
	BindNet(L)
	BindConversions(L)
	BindMime(L)
	BindComs(L)
	BindEncoding(L)
	BindMarkdown(L)
	BindLinenoise(L)
	BindOther(L)
}

// BindCarbon binds glue func
func BindCarbon(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // Carbon specific API
		"glue": glue.GetGlue,
	})
}

// BindEngine binds the engine creation.
func BindEngine(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{
		"_gin_new": gin.New,
	})
}

// BindMiddleware binds the middleware
func BindMiddleware(L *lua.State) {
	luar.Register(L, "mw", luar.Map{
		// Essentials
		"Logger":   gin.Logger,
		"Recovery": gin.Recovery,

		// Lua related stuff
		"Lua":       Lua,
		"DLR_NS":    DLR_NS,
		"DLR_RUS":   DLR_RUS,
		"DLRWS_NS":  DLRWS_NS,
		"DLRWS_RUS": DLRWS_RUS,

		// Custom sub-routers.
		"ExtRoute": (func(plan map[string]interface{}) func(*gin.Context) {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(func(*gin.Context))
			}
			return ExtRoute(newplan)
		}),
		"VHOST": (func(plan map[string]interface{}) func(*gin.Context) {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(func(*gin.Context))
			}
			return VHOST(newplan)
		}),
		"VHOST_Middleware": (func(plan map[string]interface{}) gin.HandlerFunc {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(gin.HandlerFunc)
			}
			return VHOST_Middleware(newplan)
		}),

		// To run or not to run, that is the question!
		"if_regex":       If_Regexp,
		"if_written":     If_Written,
		"if_status":      If_Status,
		"if_not_regex":   If_Not_Regexp,
		"if_not_written": If_Not_Written,
		"if_not_status":  If_Not_Status,

		// Modification stuff.
		"GZip": func() func(*gin.Context) {
			return gzip.Gzip(gzip.DefaultCompression)
		},

		// Basic
		"Echo":     EchoHTML,
		"EchoText": Echo,
	})
	luar.Register(L, "carbon", luar.Map{
		"_mw_CGI":         CGI,         // Run an CGI App!
		"_mw_CGI_Dynamic": CGI_Dynamic, // Run CGI Apps based on path!
		"_mw_combine": (func(middlewares []interface{}) func(*gin.Context) { // Combine routes, doesn't properly route like middleware or anything.
			newmiddlewares := make([]func(*gin.Context), len(middlewares))
			for k, v := range middlewares {
				newmiddlewares[k] = v.(func(*gin.Context))
			}
			return Combine(newmiddlewares)
		}),
	})
	L.DoString(glue.RouteGlue())
}

// BindStatic binds the static file server thing.
func BindStatic(L *lua.State, cfe *cache.Cache) {
	luar.Register(L, "carbon", luar.Map{
		"_staticserve": (func(path, prefix string) func(*gin.Context) {
			return staticServe.ServeCached(prefix, staticServe.PhysFS(path, prefix, true, true), cfe)
		}),
	})
}

// BindPhysFS binds the physfs library functions.
func BindPhysFS(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // PhysFS
		"_fs_mount":       physfs.Mount,
		"_fs_exists":      physfs.Exists,
		"_fs_isDir":       physfs.IsDirectory,
		"_fs_getFS":       physfs.FileSystem,
		"_fs_mkdir":       physfs.Mkdir,
		"_fs_umount":      physfs.RemoveFromSearchPath,
		"_fs_delete":      physfs.Delete,
		"_fs_setWriteDir": physfs.SetWriteDir,
		"_fs_getWriteDir": physfs.GetWriteDir,
		"_fs_list": func(name string) (fl []string, err error) {
			if physfs.Exists(name) {
				if physfs.IsDirectory(name) {
					return physfs.EnumerateFiles(name)
				}
				return nil, errors.New("readdirent: not a directory")
			}
			return nil, errors.New("open " + name + ": no such file or directory")
		},
		"_fs_readfile": func(name string) (string, error) {

			file, err := physfs.Open(name)
			if err != nil {
				return "", err
			}
			buf := bytes.NewBuffer(nil)
			io.Copy(buf, file)
			file.Close()
			return string(buf.Bytes()), nil
		},
		"_fs_modtime": func(name string) (int, error) {
			mt, err := physfs.GetLastModTime(name)
			if err != nil {
				return -1, err
			}
			return int(mt.UTC().Unix()), nil
		},
		"_fs_size": func(path string) (int64, error) {
			f, err := physfs.Open(path)
			if err != nil {
				return -1, err
			}
			info, err := f.Stat()
			if err != nil {
				return -1, err
			}
			return info.Size(), nil
		},
	})
}

// BindIOEnhancements binds small functions to enhance the IO library
func BindIOEnhancements(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // Small enhancements to the io stuff.
		"_io_list": (func(path string) ([]string, error) {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return make([]string, 1), err
			} else {
				list := make([]string, len(files))
				for i := range files {
					list[i] = files[i].Name()
				}
				return list, nil
			}
		}),
		"_io_glob": filepath.Glob,
		"_io_modtime": (func(path string) (int, error) {
			info, err := os.Stat(path)
			if err != nil {
				return -1, err
			}
			return int(info.ModTime().UTC().Unix()), nil
		}),
		"_io_isDir": func(path string) bool {
			info, err := os.Stat(path)
			if err != nil {
				return false
			}
			return info.IsDir()
		},
		"_io_size": func(path string) (int64, error) {
			info, err := os.Stat(path)
			if err != nil {
				return -1, err
			}
			return info.Size(), nil
		},
	})
}

// BindOSEnhancements does the same as above, but for the OS library
func BindOSEnhancements(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // Small enhancements to the io stuff.
		"_os_exists": (func(path string) bool {
			if _, err := os.Stat(path); err == nil {
				return true
			} else {
				return false
			}
		}),
		"_os_sleep": (func(secs int64) {
			time.Sleep(time.Duration(secs) * time.Second)
		}),
		"_os_chdir":   os.Chdir,
		"_os_abspath": filepath.Abs,
		"_os_pwd":     os.Getwd,
	})
}

// BindRedis binds the redis library
func BindRedis(L *lua.State) {
	luar.Register(L, "redis", luar.Map{
		"connectTimeout": (func(host string, timeout int) (*redis.Client, error) {
			return redis.DialTimeout("tcp", host, time.Duration(timeout)*time.Second)
		}),
		"connect": (func(host string) (*redis.Client, error) {
			return redis.Dial("tcp", host)
		}),
	})
}

// BindKVStore binds the kv store for internal carbon cache or similar.
func BindKVStore(L *lua.State) { // Thread safe Key Value Store that doesn't persist.
	luar.Register(L, "kvstore", luar.Map{
		"_set": (func(k string, v interface{}) {
			kvstore.Set(k, v, -1)
		}),
		"_del": (func(k string) {
			kvstore.Delete(k)
		}),
		"_get": (func(k string) interface{} {
			res, found := kvstore.Get(k)
			if found {
				return res
			} else {
				return nil
			}
		}),
		"_inc": (func(k string, n int64) error {
			return kvstore.Increment(k, n)
		}),
		"_dec": (func(k string, n int64) error {
			return kvstore.Decrement(k, n)
		}),
	})
}

// BindThread binds state creation and stuff.
func BindThread(L *lua.State) {
	luar.Register(L, "thread", luar.Map{
		"_spawn": (func(bcode string, dobind bool, vals map[string]interface{}, buffer int) (chan interface{}, error) {
			var ch chan interface{}
			if buffer == -1 {
				ch = make(chan interface{})
			} else {
				ch = make(chan interface{}, buffer)
			}

			L := luar.Init()
			Bind(L, webroot)
			err := L.DoString(glue.MainGlue())
			if err != nil {
				panic(err)
			}

			luar.Register(L, "", luar.Map{
				"threadcom": ch,
			})

			if dobind {
				luar.Register(L, "", vals)
			}

			if L.LoadBuffer(bcode, len(bcode), "thread") != 0 {
				return make(chan interface{}), errors.New(L.ToString(-1))
			}

			scheduler.Add(func() {
				if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
					fmt.Println("thread error: " + L.ToString(-1))
				}
			})
			return ch, nil
		}),
	})
}

// BindComs binds the com.* funcs.
func BindComs(L *lua.State) {
	luar.Register(L, "com", luar.Map{
		"create": (func() chan interface{} {
			return make(chan interface{})
		}),
		"createBuffered": (func(buffer int) chan interface{} {
			return make(chan interface{}, buffer)
		}),
		"receive": (func(c chan interface{}) interface{} {
			return <-c
		}),
		"try_receive": (func(c chan interface{}) interface{} {
			select {
			case msg := <-c:
				return msg
			default:
				return nil
			}
		}),
		"send": (func(c chan interface{}, val interface{}) bool {
			c <- val
			return true
		}),
		"try_send": (func(c chan interface{}, val interface{}) bool {
			select {
			case c <- val:
				return true
			default:
				return false
			}
		}),
		"size": (func(c chan interface{}) int {
			return len(c)
		}),
		"cap": (func(c chan interface{}) int {
			return cap(c)
		}),
		"pipe": (func(a, b chan interface{}) {
			for {
				b <- <-a
			}
		}),
		"pipe_background": (func(a, b chan interface{}) {
			scheduler.Add(func() {
				for {
					b <- <-a
				}
			})
		}),
	})
}

// BindNet binds sockets, not really that good. needs rework.
func BindNet(L *lua.State) {
	luar.Register(L, "net", luar.Map{
		"dial": net.Dial,
		"dial_tls": func(proto, addr string) (net.Conn, error) {
			config := tls.Config{InsecureSkipVerify: true} // Because I'm not gonna bother with auth.
			return tls.Dial(proto, addr, &config)
		},
		"write": (func(con interface{}, str string) {
			fmt.Fprintf(con.(net.Conn), str)
		}),
		"readline": (func(con interface{}) (string, error) {
			return bufio.NewReader(con.(net.Conn)).ReadString('\n')
		}),
		"pipe_conn": (func(con interface{}, input, output chan interface{}) {
			go func() {
				reader := bufio.NewReader(con.(net.Conn))
				for {
					line, _ := reader.ReadString('\n')
					output <- line
				}
			}()
			for {
				line := <-input
				fmt.Fprintf(con.(net.Conn), line.(string))
			}
		}),
		"pipe_conn_background": (func(con interface{}, input, output chan interface{}) {
			scheduler.Add(func() {
				reader := bufio.NewReader(con.(net.Conn))
				for {
					line, _ := reader.ReadString('\n')
					output <- line
				}
			})
			scheduler.Add(func() {
				for {
					line := <-input
					fmt.Fprintf(con.(net.Conn), line.(string))
				}
			})
		}),
	})
}

func BindMime(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{
		"_mime_byext":  mime.TypeByExtension,
		"_mime_bytype": mime.ExtensionsByType,
	})
}

// BindConversions binds helpers to convert between go types.
func BindConversions(L *lua.State) {
	luar.Register(L, "convert", luar.Map{
		"stringtocharslice": (func(x string) []byte {
			return []byte(x)
		}),
		"charslicetostring": (func(x []byte) string {
			return string(x)
		}),
	})
}

// BindEncoding binds functions to encode and decode between things.
func BindEncoding(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{
		"_enc_base64_enc": (func(str string) string {
			return base64.StdEncoding.EncodeToString([]byte(str))
		}),
		"_enc_base64_dec": (func(str string) (string, error) {
			data, err := base64.StdEncoding.DecodeString(str)
			return string(data), err
		}),
	})
}

// BindMarkdown binds a markdown renderer.
func BindMarkdown(L *lua.State) {
	luar.Register(L, "markdown", luar.Map{
		"github": (func(source string) string {
			return string(github_flavored_markdown.Markdown([]byte(source)))
		}),
	})
}

// BindLinenoise binds the linenoise library.
func BindLinenoise(L *lua.State) {
	luar.Register(L, "linenoise", luar.Map{
		"line":         linenoise.Line,
		"clear":        linenoise.Clear,
		"addHistory":   linenoise.AddHistory,
		"saveHistory":  linenoise.SaveHistory,
		"loadHistory":  linenoise.LoadHistory,
		"setMultiline": linenoise.SetMultiline,
	})
}

// BindOther binds misc things
func BindOther(L *lua.State) {
	luar.Register(L, "", luar.Map{
		"unixtime": (func() int {
			return int(time.Now().UTC().Unix())
		}),
		"regexp": regexp.Compile,
	})
	luar.Register(L, "carbon", luar.Map{
		"_syntaxhl": helpers.SyntaxHL,
	})
}

// BindContext binds the gin context.
func BindContext(L *lua.State, context *gin.Context) {
	luar.Register(L, "", luar.Map{
		"context": context,
		"req":     context.Request,

		"host":   context.Request.URL.Host,
		"path":   context.Request.URL.Path,
		"scheme": context.Request.URL.Scheme,
	})
	luar.Register(L, "carbon", luar.Map{
		"_header_set": context.Header,
		"_header_get": context.Request.Header.Get,
		"_paramfunc":  context.Param,
		"_formfunc":   context.PostForm,
		"_queryfunc":  context.Query,
	})
}
