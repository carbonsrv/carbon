package glue

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _mainglue_lua = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xdc\x57\xcd\x72\xdb\x36\x10\x3e\x9b\x4f\xc1\x20\x95\x03\xb6\x14\x63\xe7\xd0\x99\xca\xa6\x3b\x9d\xa6\xe7\x76\xda\x5c\x5a\xc9\x51\x61\x0a\xb2\x38\xa2\x40\x95\x04\x35\x4a\x55\xf9\xd9\xbb\x8b\x3f\x82\xb4\xa8\xb8\x87\x5c\x7a\xb1\x0c\xec\xee\xb7\xbb\x1f\x16\xbb\x60\x51\x66\xac\x08\x57\x72\x53\xcc\x79\x9d\xb1\x2d\x4f\x0f\x53\x72\x4b\xee\x53\x72\x59\xc8\x1b\x12\x4f\xc9\x9d\x5a\x3c\xea\xc5\xa5\x5a\xb0\xcd\xf6\x86\x1c\x83\x40\x1b\x37\x55\x6e\x6d\x97\x8d\xc8\x64\x5e\x0a\xca\xa2\xe0\xa2\xe2\xb2\xa9\x44\x48\xc9\x68\x34\xba\x7a\xb7\x27\xd1\x64\x59\x56\x1b\x26\x29\x9b\x3c\x7c\x92\x9c\x46\x51\xc0\xc5\xc2\x47\x69\xc4\x19\x9c\x5a\x56\xb9\x78\x4c\xb2\x15\xab\xa8\x2c\x45\xb3\x79\xe0\x15\x65\xf1\xf5\xb7\x16\xc7\xc6\x1f\x5c\x60\x3a\x2d\x06\x18\x02\x8a\x0b\x07\x96\x93\xc7\xba\x79\xa0\x64\x7a\x7b\x77\x79\x4f\x62\x2f\x79\x80\xba\x00\xac\x38\xb8\x68\xaa\x17\x21\x7c\x64\xe3\xbf\x7f\x18\xff\x71\x35\xfe\x6e\x9e\x3c\x8d\x01\xac\x25\xa3\xc5\xaa\x57\xbc\x78\x11\xda\xa8\xfe\xf3\xe9\xd5\xeb\xaf\x2e\xbf\xa6\xd1\x3f\xb3\xd9\x9b\x19\xb9\xb9\xbd\xfb\xfe\x70\x9c\x8e\xee\x3f\x02\x34\x99\xcd\x46\xd7\xc4\xc0\x06\x40\xbf\xa3\xeb\xf0\xd2\x78\xbf\x01\x94\x10\x0e\x42\xaf\x46\x23\x3a\xda\x8f\xf6\x91\x8e\xda\xa2\x79\x0e\xc6\xe3\x50\xb2\xc7\x70\xc3\x25\x93\xec\xa1\xe0\xe6\xa8\x60\x0f\xb6\x56\xe8\xb6\x02\x4d\x5e\x79\x9e\x79\xb1\x44\xd7\x5a\x31\x2b\x85\xe4\x42\xc2\x3a\x5f\x86\x28\x4a\xcc\x4e\x28\x57\x5c\xc0\xf6\x85\x59\xa7\x87\x23\xae\xa0\x3e\xc2\x3c\xde\x85\xb9\x08\xf3\x2d\xcb\xab\x9a\xfa\x46\x51\xb8\x28\x51\x0b\xc1\xe4\xa7\x2d\xa7\xbb\x28\x4d\x89\x2e\x0b\xe2\x10\x1d\xe6\x34\xbf\x4f\x77\x6a\x87\x17\x35\xef\xd8\xe8\xe2\x19\xb0\x91\xa5\x86\x04\x55\x67\xfd\x1c\x78\xa2\x33\xa7\x46\x07\xf8\x72\x3f\xfa\xaf\x66\xa0\xdc\x22\x2d\xb5\xc7\x80\xd9\x09\x99\x58\x84\x82\xef\x25\xf5\x77\x23\x17\x92\xd9\xf0\x88\x59\x6b\x62\x3c\x5e\x9c\x51\x8f\x97\xf5\xe9\x1c\xd5\x11\x26\xb9\xa8\x79\x25\xa9\xb1\x8d\xbd\x6c\xfb\xe9\x9e\xd4\x5f\x27\x09\x49\x67\x84\x24\x89\x67\x69\xea\x09\xb6\x63\x72\xf9\x57\x53\x42\xaf\x88\x40\x0f\xd6\x83\xf4\x58\x3a\x96\x59\x51\xd6\xdc\x45\x09\xdb\xfd\x12\xb1\x65\xac\xc3\x01\x69\x06\x2d\xc4\x96\x04\x78\xb9\x7d\x0b\xc1\x28\x2c\xc1\x36\x1c\x36\xee\x48\xeb\xcc\x1a\x0f\x68\x99\xca\xd0\x91\x94\x5b\x2e\x9c\x57\x7d\x7c\x15\xaf\x9b\x42\x9a\xc0\xec\xc9\x79\x81\xa1\x34\x85\x76\xd9\x85\x0e\x91\x1d\x3f\x5a\x4b\x1e\xde\xbd\x36\x40\x4b\xf4\x10\x4c\x27\x8f\x33\xc4\x68\xfb\x9e\x4b\x4b\xd0\x73\x2a\x5c\x4e\x36\x80\x2f\x96\xea\xe7\xb3\xfc\x6f\x09\x02\x23\x03\x59\xbe\xb4\x0c\x1c\x52\xf8\x56\x9f\x3f\x4a\x75\x93\x66\x8b\xc5\xdc\x76\xa3\x4e\x3b\x8b\x93\x24\x89\x74\xc5\x8a\x52\x9e\x6e\x63\xfe\xa6\xbe\xb2\xda\xaf\x6a\x68\xe9\x75\x0c\x72\x9e\x49\xfa\xe6\xf5\x1b\x85\x66\x2e\xac\x26\x7e\xc7\x8a\x86\xa7\x46\x23\xb7\xde\xda\x26\x87\xe2\x93\x8d\x0e\x34\x5e\x2b\xe9\x53\x7a\x35\x74\xcd\xfd\xc0\x62\xdd\xdf\x13\x9c\x77\x06\xb6\x7b\x3d\xed\x51\x0d\x23\x68\xab\xfe\x5d\xb6\x03\x1a\x34\x2d\x9b\x90\x78\xc6\xff\xef\x7c\xb6\x6c\x7c\x09\x0a\xb1\x20\xed\x14\xe8\x12\x28\x1f\x8a\x3e\x81\xfd\x0b\xeb\x6f\xf6\x09\xec\x0c\x12\xc4\x32\xf4\x9d\x1f\x1f\xcf\x33\xb2\x57\x7d\xd7\xbd\xe8\xbe\x6c\xba\x36\x33\xf8\x33\xb9\x66\x05\x67\xd5\xe9\x72\x41\xf0\x4e\x39\x88\xbc\x38\x07\x72\x92\x32\x07\x62\xa5\x83\x20\x35\x97\xe7\xcb\x16\xff\x9f\x74\xe2\xa5\x51\x17\x6a\xe2\xb5\x12\xaa\xad\x5a\xec\xb3\x27\xea\x61\x1b\xbd\x93\xd8\x56\xa6\xad\xfc\x0b\x87\x23\x6c\x20\x73\x35\xde\x52\x59\x35\xfc\xec\x8d\xc5\x69\x3c\x84\xa0\x85\x5d\x08\xbd\x32\xef\xc5\xf6\x79\x28\xf1\x71\x38\x9f\xe7\xf0\x46\xda\xf7\x72\x5d\x9b\xda\x6d\x6b\x4d\x95\x56\x5b\x6a\x43\xe9\x2a\x43\x5b\x66\xee\x51\xa0\x9e\xa2\x50\x67\xdd\x6e\x3e\x9f\x43\x24\x45\x6a\xc4\x89\x77\x22\xe6\x5d\xfb\xdb\x96\x67\xf9\x32\xcf\xc2\x0f\xec\xb1\x0e\x6c\x84\x88\x47\x71\x7e\x78\x9f\x1d\x5c\xba\xf7\x2f\x3d\xa0\x2c\xc5\x3f\xc7\x58\xa5\x69\x3e\x3d\x9c\x7d\x91\x8b\x35\x85\xb7\x78\x8c\x23\x14\x40\xf4\x2c\xb5\xa9\xb5\x41\x13\x46\xa6\x87\x55\xc5\x97\xa9\x52\x6e\xc4\x96\x65\x6b\x1c\x9e\xd1\x11\x32\x31\x49\x0e\xa8\x2b\x0d\xf0\xda\xf5\x5c\x67\x55\xbe\xc5\x89\xb8\xf0\x82\x47\x53\x2d\x20\x46\xd2\x35\x5a\x94\x99\x3a\x86\xae\xc5\xab\xf7\x3f\xff\xf8\xe1\xf7\x5f\x7e\x02\x9f\x04\x87\x05\x39\xde\x4f\xda\xfa\xa2\x1a\x04\x38\xfc\x55\x9b\x58\xb4\x16\xd6\x16\xff\x02\x68\x8b\x43\x74\x0c\x7f\xd1\x11\xf8\xb1\x5f\x06\x0b\x1e\xa6\xfa\x07\x5a\xd2\xbb\xab\x2b\x27\x41\x3d\x14\xa9\x5f\x90\x11\x09\xaf\xe4\xb7\x2a\x8e\xa0\xfb\x5d\x01\x5a\x84\x04\xae\x96\xd0\x5b\x14\xa6\xb0\xd9\x6b\xed\xad\x3a\xaa\x04\x9d\xef\x81\xd6\xc8\x54\x21\xbe\xcb\x71\x33\xd1\x0f\xfc\xf0\x29\x0d\xa1\x59\x9c\x86\xf2\x3e\x02\xcc\x91\xb5\x72\xf7\x3c\x56\x1e\xf4\x89\x69\xf1\x5e\x26\xef\x61\x8f\x7a\xbc\x20\x49\x62\x07\x9d\x15\x5e\xd5\xf8\x81\x5b\x17\x79\xc6\xdd\xeb\xa6\x7f\x6a\xf8\x19\x6d\xab\x14\x87\x00\xfc\xd7\x0b\x73\x09\x01\xcc\x51\x0d\x6d\xa8\x0b\x45\xd9\x98\xeb\xb7\x54\x49\x3f\xbb\x75\xba\x31\x76\x1a\xf5\xb2\x2d\xb7\x7f\x03\x00\x00\xff\xff\x66\x58\x64\xd6\x2e\x10\x00\x00")

func mainglue_lua_bytes() ([]byte, error) {
	return bindata_read(
		_mainglue_lua,
		"MainGlue.lua",
	)
}

func mainglue_lua() (*asset, error) {
	bytes, err := mainglue_lua_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "MainGlue.lua", size: 4142, mode: os.FileMode(420), modTime: time.Unix(1432575618, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _routeglue_lua = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x52\xcb\x8e\xd4\x40\x0c\x3c\x27\x5f\x61\xe5\x94\x48\xd9\xd5\x8a\x23\xd2\xdc\xe0\x06\x1c\x16\x38\x21\x84\xf2\x70\x48\x4b\x3d\xdd\x91\xdb\xd9\x30\x42\xf0\xed\xd8\x9d\xe7\xa2\x19\x24\xf6\x62\x4d\xda\x55\xae\x72\x8d\xef\xee\xe0\xbd\x69\x5b\x8b\x53\x45\x08\x3d\xda\x01\x29\x94\xe0\x3c\x03\x61\x65\xed\x05\xc6\x80\xdd\x68\x61\x32\xdc\x43\xe5\x2e\xdc\x1b\xf7\x1d\xea\x91\x81\x7b\x84\x80\xf4\x84\x04\xc6\x19\x86\xd0\x90\x19\x38\xed\x46\xd7\xb0\xf1\x0e\xce\xd3\xbd\xc3\x29\xef\x5c\x09\xb5\x71\xad\xd0\x74\x30\x4e\x81\x2b\xc6\x22\x4d\xac\x6f\x2a\x0b\x8d\x6f\x11\x4e\x90\x65\x69\x62\x3a\xe0\xcb\x80\xc2\x28\xe0\x24\x4f\xeb\xa4\x4c\xa5\x5c\x9a\x24\x0b\x36\x30\xc9\xb0\xfb\x76\x3c\x0f\x8a\x4d\x13\xb4\x01\xff\x26\xcf\xa0\x8d\xaa\x2e\x90\x48\xd8\xd6\x57\xed\xdc\xcc\x75\x9e\xd0\x55\x58\x17\xd6\xfe\x02\xbf\x29\x15\xb5\x14\x20\x60\x4f\xb9\xd4\xf8\xe8\xda\x74\x2e\xf3\x52\xb4\xfe\x90\x7e\xdc\x6b\x5d\x7b\x9d\xbf\xba\x5d\x83\x99\x3d\x73\x55\x5b\xdc\x2c\x27\xb4\x5a\x96\x24\xdf\xbc\x7b\xfc\xf6\xe1\x63\x74\x5c\x02\xd3\x88\x7b\xa8\x47\x57\x37\x29\x5d\x25\x88\x12\x7e\x7e\xc9\x42\xf6\xf5\x94\x3d\x65\xbf\x76\xdf\x33\xf7\x45\x96\x1e\x3f\xff\xbf\xa7\x9d\xf3\x2f\x53\x5a\xc4\x92\x32\x7f\x9f\xc0\x19\xbb\x7a\x38\x06\x1f\x51\x84\x3c\x92\x93\xcc\xf5\xeb\x78\x7d\xd8\xf4\x7e\x51\x22\x0c\xc3\x76\x72\xfa\x21\x76\xd8\xbb\xf1\x5c\x23\xe5\xb1\x09\x9e\xe0\xd5\xc3\xc3\x7e\x85\xf1\x38\xae\x9d\xd2\xa2\x27\x02\x6f\x55\x40\xd9\x25\x2c\xa7\x74\xbc\xc4\x7d\xc0\xf3\x10\x6f\xf2\x5f\x93\x6c\x20\x7e\x8a\x62\xfb\x53\x76\xac\x2e\xf3\x09\x7f\xf0\x82\x67\x7f\xbc\xe1\x25\x8a\x6b\x01\x44\x0e\x4b\xd9\x42\x78\xae\x1f\xfb\xf9\xd5\x2c\x0a\x91\x91\x6e\x11\xe7\xfe\x09\x00\x00\xff\xff\x04\x22\x86\x71\x29\x04\x00\x00")

func routeglue_lua_bytes() ([]byte, error) {
	return bindata_read(
		_routeglue_lua,
		"RouteGlue.lua",
	)
}

func routeglue_lua() (*asset, error) {
	bytes, err := routeglue_lua_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "RouteGlue.lua", size: 1065, mode: os.FileMode(420), modTime: time.Unix(1432300032, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"MainGlue.lua": mainglue_lua,
	"RouteGlue.lua": routeglue_lua,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"MainGlue.lua": &_bintree_t{mainglue_lua, map[string]*_bintree_t{
	}},
	"RouteGlue.lua": &_bintree_t{routeglue_lua, map[string]*_bintree_t{
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

