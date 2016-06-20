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
	"io"
	"io/ioutil"
	"math/rand"
	"mime"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	"os/exec"

	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/GeertJohan/go.linenoise"
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/carbonsrv/carbon/modules/helpers"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/carbonsrv/carbon/modules/static"
	"github.com/fzzy/radix/redis"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/pmylund/go-cache"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/vifino/contrib/gzip"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"

	// DB stuff
	"database/sql"
	_ "github.com/cznic/ql/driver"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/vifino/go-yql"
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
	BindSQL(L)
	BindKVStore(L)
	BindPhysFS(L)
	BindIOEnhancements(L)
	BindOSEnhancements(L)
	BindThread(L)
	BindNet(L)
	BindConversions(L)
	BindMime(L)
	BindComs(L)
	BindPipe(L)
	BindExec(L)
	BindEncoding(L)
	BindMarkdown(L)
	BindLinenoise(L)
	BindTermbox(L)
	BindOther(L)
	BindRandomString(L)
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
		"_physfs_mount":       physfs.Mount,
		"_physfs_exists":      physfs.Exists,
		"_physfs_isDir":       physfs.IsDirectory,
		"_physfs_getFS":       physfs.FileSystem,
		"_physfs_mkdir":       physfs.Mkdir,
		"_physfs_umount":      physfs.RemoveFromSearchPath,
		"_physfs_delete":      physfs.Delete,
		"_physfs_setWriteDir": physfs.SetWriteDir,
		"_physfs_getWriteDir": physfs.GetWriteDir,
		"_physfs_list": func(name string) (fl []string, err error) {
			if physfs.Exists(name) {
				if physfs.IsDirectory(name) {
					return physfs.EnumerateFiles(name)
				}
				return nil, errors.New("readdirent: not a directory")
			}
			return nil, errors.New("open " + name + ": no such file or directory")
		},
		"_physfs_readfile": func(file string) (string, error) {
			if physfs.Exists(file) {
				f, err := physfs.Open(file)
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
					if err.Error() == "EOF" { // Hack. Sometimes, things just don't work. No idea why.
						f2, _ := physfs.Open(file)
						buf := bytes.NewBuffer(nil)
						io.Copy(buf, f2)
						f.Close()
						return string(buf.Bytes()), nil
					}
					return "", err
				}
				return string(buf), err
			}
			return "", errors.New(file + ": No such file or directory")
		},
		"_physfs_modtime": func(name string) (int, error) {
			mt, err := physfs.GetLastModTime(name)
			if err != nil {
				return -1, err
			}
			return int(mt.UTC().Unix()), nil
		},
		"_physfs_size": func(path string) (int64, error) {
			f, err := physfs.Open(path)
			defer f.Close()
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
			}
			list := make([]string, len(files))
			for i := range files {
				list[i] = files[i].Name()
			}
			return list, nil
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
			}
			return false
		}),
		"_os_sleep": (func(milliseconds int64) {
			time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		}),
		"_os_chdir":   os.Chdir,
		"_os_abspath": filepath.Abs,
		"_os_pwd":     os.Getwd,
	})
}

// BindRedis binds the redis library
func BindRedis(L *lua.State) {
	luar.Register(L, "redis", luar.Map{
		"connectTimeout": (func(host string, timeout float64) (*redis.Client, error) {
			return redis.DialTimeout("tcp", host, time.Duration(timeout)*time.Second)
		}),
		"connect": (func(host string) (*redis.Client, error) {
			return redis.Dial("tcp", host)
		}),
	})
}

type sql_rows struct {
	Len  int
	Vals []map[string]interface{}
}

// BindSQL binds the database library
func BindSQL(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{
		"_sql_drivers": sql.Drivers,
		"_sql_open":    sql.Open,

		"_sql_rows": (func(rows *sql.Rows) (sql_rows, error) {
			var res = make([]map[string]interface{}, 0)
			rowno := 0
			names, err := rows.Columns()
			if err != nil {
				return sql_rows{}, err
			}
			numnames := len(names)

			for rows.Next() {
				//fmt.Println("Generating interfaces...")
				var elems = make([]interface{}, numnames)
				for i := range elems {
					var ii interface{}
					elems[i] = &ii
				}

				//fmt.Println("Scanning...")
				err := rows.Scan(elems...)
				if err != nil {
					return sql_rows{}, err
				}

				//fmt.Println("Storing in map...")
				rowtmp := make(map[string]interface{})
				for i, name := range names {
					var elm_val = *(elems[i].(*interface{}))
					var elm_type = reflect.TypeOf(elm_val).String()
					//fmt.Println(elm_type)
					if elm_type == "[]uint8" {
						rowtmp[name] = string(elm_val.([]uint8))
					} else {
						rowtmp[name] = elm_val
					}
				}

				//fmt.Println("Adding to slice...")
				n := len(res)
				newres := make([]map[string]interface{}, n+1)
				//fmt.Printf("Len: %d\n", n)
				copy(newres, res)
				newres[n] = rowtmp
				res = newres

				//fmt.Printf("Row %d: %#v\n", rowno, rowtmp)
				rowno += 1
			}

			final := sql_rows{
				Len:  len(res),
				Vals: res,
			}

			err = rows.Err()
			if err != nil {
				return final, err
			} else {
				return final, nil
			}
		}),
		"_sql_getrow": (func(rows sql_rows, n int) (map[string]interface{}, error) {
			if n <= rows.Len && n > 0 {
				return rows.Vals[n-1], nil
			}
			return make(map[string]interface{}), errors.New("out of bounds")
		}),
	})
}

// BindKVStore binds the kv store for internal carbon cache or similar.
func BindKVStore(L *lua.State) { // Thread safe Key Value Store that doesn't persist.
	luar.Register(L, "kvstore", luar.Map{
		"_set": (func(k string, v interface{}) {
			kvstore.Set(k, v, -1)
		}),
		"_set_timeout": (func(k string, v interface{}, timeout int64) {
			kvstore.Set(k, v, time.Duration(timeout)*time.Millisecond)
		}),
		"_del": (func(k string) {
			kvstore.Delete(k)
		}),
		"_get": (func(k string) interface{} {
			res, found := kvstore.Get(k)
			if found {
				return res
			}
			return nil
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
				L.Close()
				close(ch)
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
		"close": (func(c chan interface{}) {
			close(c)
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

// Pipe is the struct containing the reader and writer of a io.Pipe
type Pipe struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

// Read from the pipe
func (p *Pipe) Read(bytes int) (error, error) {
	b := make([]byte, bytes)
	_, err := p.reader.Read(b)
	if err != nil {
		return nil, err
	}
	return errors.New(string(b)), nil
}

// Write to the pipe
func (p *Pipe) Write(text string) error {
	_, err := p.writer.Write([]byte(text))
	return err
}

// Close it, dammit!
func (p *Pipe) Close() error {
	return p.reader.Close()
}

// BindPipe binds pipe helpers
func BindPipe(L *lua.State) {
	luar.Register(L, "io", luar.Map{
		"pipe": func() Pipe {
			r, w := io.Pipe()
			return Pipe{
				writer: w,
				reader: r,
			}
		},
	})
}

// Command is the struct containing the reader and writer of a io.Pipe
type Command struct {
	stdout io.ReadCloser
	stderr io.ReadCloser
	stdin  io.WriteCloser
	cmd    *exec.Cmd
}

// Read_Stdout reads stdout from the command
func (c *Command) Read_Stdout(bytes int) (error, error) {
	b := make([]byte, bytes)
	_, err := c.stdout.Read(b)
	if err != nil {
		return nil, err
	}
	return errors.New(string(b)), nil
}

// Read_Stderr reads stderr from the command
func (c *Command) Read_Stderr(bytes int) (error, error) {
	b := make([]byte, bytes)
	_, err := c.stderr.Read(b)
	if err != nil {
		return nil, err
	}
	return errors.New(string(b)), nil
}

// Write_Stdin writes to the command's stdin
func (c *Command) Write_Stdin(text string) error {
	_, err := c.stdin.Write([]byte(text))
	return err
}

// Close it, dammit!
func (c *Command) Close() {
	c.stdin.Close()
	c.stdout.Close()
	c.stderr.Close()
}

// Kill misbehaving programs.
func (c *Command) Kill() error {
	e := c.cmd.Process.Kill()
	c.Close()
	return e
}

// Exited checks if the process exited already
func (c *Command) State() *os.ProcessState {
	return c.cmd.ProcessState
}

// BindExec binds exec.exec to call shit.
func BindExec(L *lua.State) {
	luar.Register(L, "exec", luar.Map{
		"exec": func(cmd string, args ...string) (Command, error) {
			command := exec.Command(cmd, args...)
			stdin, err := command.StdinPipe()
			if err != nil {
				return Command{}, err
			}
			stdout, err := command.StdoutPipe()
			if err != nil {
				return Command{}, err
			}
			stderr, err := command.StderrPipe()
			if err != nil {
				return Command{}, err
			}
			c := Command{
				stdin:  stdin,
				stdout: stdout,
				stderr: stderr,
				cmd:    command,
			}
			if err := command.Start(); err != nil {
				return Command{}, err
			}
			/*scheduler.Add(func() {
				command.Wait()
			})*/
			return c, nil
		},
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
		"listen": net.Listen,
		"write": (func(con interface{}, str string) {
			fmt.Fprintf(con.(net.Conn), str)
		}),
		"readline": (func(con interface{}) (string, error) {
			return bufio.NewReader(con.(net.Conn)).ReadString('\n')
		}),
		"read": (func(con interface{}, count int) (string, error) {
			b := make([]byte, count)
			_, err := con.(net.Conn).Read(b)
			if err != nil {
				return "", err
			}
			return string(b), nil
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

// BindMime binds mime lookup functions
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
	lua_complfn := func(L *lua.State) int {
		if L.LTypename(1) != "function" {
			// TODO: error
		}

		L.PushValue(1)
		compfn_ref := L.Ref(lua.LUA_REGISTRYINDEX)

		str_slice_type := reflect.TypeOf([]string{})

		linenoise.SetCompletionHandler(func(in string) []string {
			L.RawGeti(lua.LUA_REGISTRYINDEX, compfn_ref)
			L.PushString(in)
			if L.Pcall(1, 1, 0) != 0 { // call compfn with (1) in string, 1 return value
				// TODO: error maybe?
				return []string{}
			}
			if L.LTypename(-1) != "table" {
				// TODO: error or something
				return []string{}
			}
			str_slice := luar.CopyTableToSlice(L, str_slice_type, -1)
			return str_slice.([]string)
		})
		return 1
	}
	//L.Register("linenoise.SetCompletionHandler", lua_complfn)
	L.GetGlobal("linenoise")
	L.PushString("setCompletionHandler")
	L.PushGoFunction(lua_complfn)
	L.SetTable(-3)
}

// BindTermbox binds the termbox-go library.
func BindTermbox(L *lua.State) {
	luar.Register(L, "termbox", luar.Map{
		// Functions
		"CellBuffer": termbox.CellBuffer,
		"Clear":      termbox.Clear,
		"Close":      termbox.Close,
		"Flush":      termbox.Flush,
		"HideCursor": termbox.HideCursor,
		"Init":       termbox.Init,
		"Interrupt":  termbox.Interrupt,
		"SetCell": func(x, y int, ch string, fg, bg termbox.Attribute) {
			termbox.SetCell(x, y, []rune(ch)[0], fg, bg)
		},
		"SetCursor":     termbox.SetCursor,
		"Size":          termbox.Size,
		"Sync":          termbox.Sync,
		"SetInputMode":  termbox.SetInputMode,
		"SetOutputMode": termbox.SetOutputMode,
		"PollEventRaw": func() map[string]interface{} {
			e := termbox.PollEvent()
			return map[string]interface{}{
				"Type":   e.Type,
				"Mod":    e.Mod,
				"Key":    e.Key,
				"Ch":     string(e.Ch),
				"Width":  e.Width,
				"Height": e.Height,
				"Err":    e.Err,
				"MouseX": e.MouseX,
				"MouseY": e.MouseY,
				"N":      e.N,
			}
		},
		"PollRawEvent": termbox.PollRawEvent,
		"ParseEvent":   termbox.ParseEvent,

		"TBPrint": func(x, y int, fg, bg termbox.Attribute, msg string) {
			for _, c := range msg {
				termbox.SetCell(x, y, c, fg, bg)
				x += runewidth.RuneWidth(c)
			}
		},

		// Constants:
		// Colors
		"ColorDefault": termbox.ColorDefault,
		"ColorBlack":   termbox.ColorBlack,
		"ColorRed":     termbox.ColorRed,
		"ColorGreen":   termbox.ColorGreen,
		"ColorYellow":  termbox.ColorYellow,
		"ColorBlue":    termbox.ColorBlue,
		"ColorMagenta": termbox.ColorMagenta,
		"ColorCyan":    termbox.ColorCyan,
		"ColorWhite":   termbox.ColorWhite,
		// Attributes
		"AttrBold":      termbox.AttrBold,
		"AttrUnderline": termbox.AttrUnderline,
		"AttrReverse":   termbox.AttrReverse,
		// Events
		"EventKey":       termbox.EventKey,
		"EventResize":    termbox.EventResize,
		"EventMouse":     termbox.EventMouse,
		"EventError":     termbox.EventError,
		"EventInterrupt": termbox.EventInterrupt,
		"EventRaw":       termbox.EventRaw,
		"EventNone":      termbox.EventNone,

		// Input Modes
		"InputEsc":     termbox.InputEsc,
		"InputAlt":     termbox.InputAlt,
		"InputMouse":   termbox.InputMouse,
		"InputCurrent": termbox.InputCurrent,

		// Output modes
		"OutputCurrent":   termbox.OutputCurrent,
		"OutputNormal":    termbox.OutputNormal,
		"Output256":       termbox.Output256,
		"Output216":       termbox.Output216,
		"OutputGrayscale": termbox.OutputGrayscale,

		// Keys
		"KeyEsc":        termbox.KeyEsc,
		"KeySpace":      termbox.KeySpace,
		"KeyTab":        termbox.KeyTab,
		"KeyEnter":      termbox.KeyEnter,
		"KeyF1":         termbox.KeyF1,
		"KeyF2":         termbox.KeyF2,
		"KeyF3":         termbox.KeyF3,
		"KeyF4":         termbox.KeyF4,
		"KeyF5":         termbox.KeyF5,
		"KeyF6":         termbox.KeyF6,
		"KeyF7":         termbox.KeyF7,
		"KeyF8":         termbox.KeyF8,
		"KeyF9":         termbox.KeyF9,
		"KeyF10":        termbox.KeyF10,
		"KeyF11":        termbox.KeyF11,
		"KeyF12":        termbox.KeyF12,
		"KeyInsert":     termbox.KeyInsert,
		"KeyDelete":     termbox.KeyDelete,
		"KeyHome":       termbox.KeyHome,
		"KeyEnd":        termbox.KeyEnd,
		"KeyPgup":       termbox.KeyPgup,
		"KeyPgdn":       termbox.KeyPgdn,
		"KeyArrowUp":    termbox.KeyArrowUp,
		"KeyArrowDown":  termbox.KeyArrowDown,
		"KeyArrowLeft":  termbox.KeyArrowLeft,
		"KeyArrowRight": termbox.KeyArrowRight,

		// More Keys
		"KeyCtrlTilde":         termbox.KeyCtrlTilde,
		"KeyCtrlSpace":         termbox.KeyCtrlSpace,
		"KeyBackspace":         termbox.KeyBackspace,
		"KeyBackspace2":        termbox.KeyBackspace2,
		"KeyCtrlSlash":         termbox.KeyCtrlSlash,
		"KeyCtrlBackslash":     termbox.KeyCtrlBackslash,
		"KeyCtrlLsqBracket":    termbox.KeyCtrlLsqBracket,
		"KeyKeyCtrlRsqBracket": termbox.KeyCtrlRsqBracket,
		"KeyCtrlUnderscore":    termbox.KeyCtrlUnderscore,

		// Ctrl everything
		"KeyCtrl2": termbox.KeyCtrl2,
		"KeyCtrl3": termbox.KeyCtrl3,
		"KeyCtrl4": termbox.KeyCtrl4,
		"KeyCtrl5": termbox.KeyCtrl5,
		"KeyCtrl6": termbox.KeyCtrl6,
		"KeyCtrl7": termbox.KeyCtrl7,
		"KeyCtrl8": termbox.KeyCtrl8,
		"KeyCtrlA": termbox.KeyCtrlA,
		"KeyCtrlB": termbox.KeyCtrlB,
		"KeyCtrlC": termbox.KeyCtrlC,
		"KeyCtrlD": termbox.KeyCtrlD,
		"KeyCtrlE": termbox.KeyCtrlE,
		"KeyCtrlF": termbox.KeyCtrlF,
		"KeyCtrlG": termbox.KeyCtrlG,
		"KeyCtrlH": termbox.KeyCtrlH,
		"KeyCtrlI": termbox.KeyCtrlI,
		"KeyCtrlJ": termbox.KeyCtrlJ,
		"KeyCtrlK": termbox.KeyCtrlK,
		"KeyCtrlL": termbox.KeyCtrlL,
		"KeyCtrlM": termbox.KeyCtrlM,
		"KeyCtrlN": termbox.KeyCtrlN,
		"KeyCtrlO": termbox.KeyCtrlO,
		"KeyCtrlP": termbox.KeyCtrlP,
		"KeyCtrlQ": termbox.KeyCtrlQ,
		"KeyCtrlR": termbox.KeyCtrlR,
		"KeyCtrlS": termbox.KeyCtrlS,
		"KeyCtrlT": termbox.KeyCtrlT,
		"KeyCtrlU": termbox.KeyCtrlU,
		"KeyCtrlV": termbox.KeyCtrlV,
		"KeyCtrlW": termbox.KeyCtrlW,
		"KeyCtrlX": termbox.KeyCtrlX,
		"KeyCtrlY": termbox.KeyCtrlY,
		"KeyCtrlZ": termbox.KeyCtrlZ,

		// Modifiers
		"ModAlt":    termbox.ModAlt,
		"ModMotion": termbox.ModMotion,

		// Mouse
		"MouseLeft":      termbox.MouseLeft,
		"MouseMiddle":    termbox.MouseMiddle,
		"MouseRight":     termbox.MouseRight,
		"MouseRelease":   termbox.MouseRelease,
		"MouseWheelUp":   termbox.MouseWheelUp,
		"MouseWheelDown": termbox.MouseWheelDown,
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func BindRandomString(L *lua.State) {
	var random_src = rand.NewSource(time.Now().UnixNano())
	luar.Register(L, "carbon", luar.Map{
		"randomstring": func(n int) string {
			b := make([]byte, n)
			// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
			for i, cache, remain := n-1, random_src.Int63(), letterIdxMax; i >= 0; {
				if remain == 0 {
					cache, remain = random_src.Int63(), letterIdxMax
				}
				if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
					b[i] = letterBytes[idx]
					i--
				}
				cache >>= letterIdxBits
				remain--
			}

			return string(b)
		},
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
