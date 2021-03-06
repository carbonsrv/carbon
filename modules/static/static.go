package staticServe

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
)

// ServeFileSystem is a small wrapper
type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	origfs  http.FileSystem
	root    string
	prefix  string
	indexes bool
	physfs  bool
}

// OwnFS returns the fs wrapper using an existing filesystem
func OwnFS(fs http.FileSystem, root, prefix string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: fs,
		origfs:     fs,
		root:       "/root/" + root,
		prefix:     prefix,
		indexes:    indexes,
		physfs:     false,
	}
}

// PhysFS returns the fs wrapper using physfs
func PhysFS(root, prefix string, indexes bool, alreadyinitialized bool) *localFileSystem {
	if !alreadyinitialized {
		root, err := filepath.Abs(root)
		fmt.Println(root)
		if err != nil {
			panic(err)
		}
		err = physfs.Init()
		if err != nil {
			panic(err)
		}
		defer physfs.Deinit()
		err = physfs.Mount(root, "/", true)
		if err != nil {
			panic(err)
		}
	}
	fs := physfs.FileSystem()
	return &localFileSystem{
		FileSystem: fs,
		origfs:     fs,
		root:       root,
		prefix:     prefix,
		indexes:    indexes,
		physfs:     true,
	}
}

// Open is the function the wrapper is about
func (l *localFileSystem) Open(name string) (http.File, error) {
	//if !l.physfs {
	newname := name
	if p := strings.TrimPrefix(name, l.prefix); len(p) <= len(name) {
		newname = path.Join(l.root, p)
	}
	f, err := l.origfs.Open(newname)
	if err != nil {
		return nil, err
	}
	if l.indexes {
		return f, err
	} else {
		return neuteredReaddirFile{f}, nil
	}
	/*} else {
		return physfs.Open(name)
	}*/
}

// Exists is also a function for the wrapper
func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) <= len(filepath) {
		p = path.Join(l.root, p)
		/*if !l.physfs {
			return existsFile(l, p)
		} else {*/
		fmt.Println("Exists: " + p)
		return physfs.Exists(p)
		//}
	}
	return false
}

func cachedFileExists(fs *localFileSystem, cfe *cache.Cache, prefix string, path string) bool {
	data_tmp, found := cfe.Get(path)
	if found == false {
		exists := fs.Exists(prefix, path)
		cfe.Set(path, exists, cache.DefaultExpiration)
		return exists
	} else {
		return data_tmp.(bool)
	}
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

// ServeCached takes a fs wrapper and serves stuff.
func ServeCached(prefix string, fs ServeFileSystem, cfe *cache.Cache) gin.HandlerFunc {
	//cfe := cache.New(5*time.Minute, 30*time.Second) // File-Exists Cache
	var fileserver http.Handler

	if prefix != "" && prefix != "/" {
		fileserver = http.StripPrefix(prefix, http.FileServer(fs))
	} else {
		fileserver = http.FileServer(fs)
	}

	return func(c *gin.Context) {
		//if cachedFileExists(fs, cfe, prefix, c.Request.URL.Path) {
		fileserver.ServeHTTP(c.Writer, c.Request)
		/*} else {
			c.Abort()
		}*/
	}
}
