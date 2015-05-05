package staticServe

import (
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	fs      http.FileSystem
	root    string
	indexes bool
	physfs  bool
}

func existsFile(lfs *localFileSystem, name string) bool {
	if !lfs.physfs {
		_, err := os.Stat(name)
		return !os.IsNotExist(err)
	} else {
		return physfs.Exists(name)
	}
}

func LocalFile(root string, indexes bool) *localFileSystem {
	root, err := filepath.Abs(root)
	fmt.Println(root)
	if err != nil {
		panic(err)
	}

	fs := http.Dir(root)
	return &localFileSystem{
		fs,
		root,
		indexes,
		false,
	}
}

func OwnFS(fs http.FileSystem, root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		fs,
		root,
		indexes,
		false,
	}
}

func PhysFS(root string, indexes bool, alreadyinitialized bool) *localFileSystem {
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
		fs,
		root,
		indexes,
		true,
	}
}

func (l *localFileSystem) Open(name string) (http.File, error) {
	if !l.physfs {
		f, err := l.fs.Open(name)
		if err != nil {
			return nil, err
		}
		if l.indexes {
			return f, err
		} else {
			return neuteredReaddirFile{f}, nil
		}
	} else {
		return physfs.Open(name)
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		p = path.Join(l.root, p)
		if !l.physfs {
			return existsFile(l, p)
		} else {
			return physfs.Exists(p)
		}
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

func ServeCached(prefix string, fs *localFileSystem, cfe *cache.Cache) gin.HandlerFunc {
	//cfe := cache.New(5*time.Minute, 30*time.Second) // File-Exists Cache
	var fileserver http.Handler

	if prefix != "" {
		fileserver = http.StripPrefix(prefix, http.FileServer(fs.fs))
	} else {
		fileserver = http.FileServer(fs.fs)
	}

	return func(c *gin.Context) {
		if cachedFileExists(fs, cfe, prefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
		} else {
			c.Next()
		}
	}
}
