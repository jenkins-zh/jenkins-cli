// Code generated for package i18n by go-bindata DO NOT EDIT. (@generated)
// sources:
// zh_CN/LC_MESSAGES/k8s.mo
// zh_CN/LC_MESSAGES/k8s.po
package i18n

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _zh_cnLc_messagesK8sMo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x90\x3f\x6f\x14\x31\x14\xc4\x27\x11\x34\x2e\xa9\x29\xde\x07\xc0\x47\x12\x81\x14\xf9\x38\x9a\x90\xa0\xa0\x04\x16\x74\x20\x3a\x64\x76\x1f\xb7\x0e\x3e\x7b\x65\x7b\x09\xa1\x42\xe2\x4f\x41\x95\x06\x09\xa5\x02\x09\x0a\x52\x84\x32\xdb\xf0\x55\xa8\x76\xb9\x9e\x4f\x80\xc2\x72\xb9\xa9\x7e\xa3\xf7\x46\x1a\xcd\xaf\x4b\x17\x3e\x02\xc0\x32\x80\xcb\x00\xae\x00\xb8\x08\xe0\x06\x7a\x65\x00\x6e\x02\xb8\x0f\xe0\x0f\x80\x43\x00\xd7\x01\xbc\x59\x02\x96\xb0\xd0\xf2\x1c\xf6\x72\x6b\xc8\x44\xd2\x94\xbc\xb7\xb4\x5f\x9a\xbc\xa4\xdc\xd7\xb6\xa0\x92\x6d\x45\x07\xbe\xa6\x7d\x93\xca\x33\x08\x34\xad\x6d\x32\x95\x65\xba\xc3\xee\xb9\x71\x11\x59\xf0\x7b\x9c\x27\xb9\x5d\xc8\x47\x1c\xa2\xf1\x4e\x91\xc8\xee\xc9\x07\xfc\xc2\x9c\x39\x79\x4b\x27\x56\x24\x76\xb4\x9b\xd4\x7a\xc2\x72\xcc\x7a\xaa\x48\xec\x6e\xef\x6e\x2e\x12\xab\x83\x15\xb1\xe1\x5d\x62\x97\xe4\xf8\xa0\x62\x45\x89\x5f\xa6\xab\x95\xd5\xc6\x0d\x29\x2f\x75\x88\x9c\x46\x0f\xc7\x5b\x72\x7d\xf1\x17\xb4\x8b\xcf\x38\xc8\x4d\x97\xfb\xc2\xb8\x89\xa2\xf5\xa7\x26\x89\xc7\xf2\x36\x3b\x0e\x3a\xf9\xa0\x28\xf3\x5c\x98\x44\x6b\x83\xb5\xc1\x35\xb1\xa3\xe3\xff\x98\xed\xaf\x22\xb3\x75\xd0\x56\x6e\xf9\x30\x8d\x8a\x5c\xf5\xcf\xc6\xd1\xea\x90\x7a\x1c\xad\x0c\xcf\x9b\x2b\x7a\x55\x3e\xd9\xb8\x2b\xfa\xc5\x7e\x7f\xfa\xd1\x36\xaf\xdb\xe6\xb8\x6b\x4e\xba\x0f\xdf\xdb\x9f\x9f\x67\x27\x5f\x66\x87\xef\xbb\xaf\x47\x6d\x73\x3c\x9f\x87\x66\x47\x6f\xbb\xd3\x6f\xdd\xbb\x53\xfc\x0d\x00\x00\xff\xff\x23\x70\x78\x8d\xb9\x01\x00\x00")

func zh_cnLc_messagesK8sMoBytes() ([]byte, error) {
	return bindataRead(
		_zh_cnLc_messagesK8sMo,
		"zh_CN/LC_MESSAGES/k8s.mo",
	)
}

func zh_cnLc_messagesK8sMo() (*asset, error) {
	bytes, err := zh_cnLc_messagesK8sMoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "zh_CN/LC_MESSAGES/k8s.mo", size: 441, mode: os.FileMode(420), modTime: time.Unix(1574260489, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _zh_cnLc_messagesK8sPo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x91\x3f\x6f\xd4\x40\x10\xc5\x7b\x7f\x8a\x91\x53\xaf\x9d\x9c\x82\x14\xed\xe9\xaa\x23\x41\x41\x04\x22\x64\x10\x05\x52\xb4\x59\x0f\xf6\xe6\xf6\x8f\xb5\x3b\x4e\x08\x15\x05\x50\x50\x51\xa2\x54\x20\x41\x41\x8a\x50\x9e\x1b\xbe\x8d\xcd\x7d\x0d\x64\x9f\xe1\x74\x74\x6f\x7f\xfb\xe6\x69\xe7\xad\x09\x85\xca\x21\x8e\x23\x13\x8a\x40\xbe\x57\xf1\xa9\x77\x17\x28\x89\x1d\xe7\xec\x39\xfa\xa0\x9c\xe5\xf0\xd2\xf6\x17\x4f\x32\x36\xf7\x28\x48\x39\xcb\xee\x0b\xc2\x7f\x9c\x3d\xc5\x4b\x15\xfe\xc3\x8f\x84\x2d\x6a\x51\x20\xcb\x50\x98\x91\x9d\x1c\x9f\x1c\x6e\x52\xf7\x92\xdd\x81\xce\x9d\x25\xb4\xc4\xb2\xeb\x0a\x39\x10\xbe\xa6\xb4\xd2\x42\xd9\x29\xc8\x52\xf8\x80\x34\x7b\x96\x1d\xb1\x83\x6d\xaf\x17\x36\xbc\x42\xcf\x0e\xad\x74\xb9\xb2\x05\x87\x83\x73\x45\x83\xe7\x05\x7b\x80\x16\xbd\x20\xe7\x39\x9c\x3a\xcc\x15\xc1\x24\x99\x24\xfb\xe3\xbb\xc2\x38\xae\xd7\x8e\xf5\x12\xba\xf6\x42\xb3\x23\xe7\x4d\xe0\x60\xab\xe1\x18\x66\x7b\x53\x58\xcb\xd9\xee\x74\x6b\x2b\x0e\x6f\xca\xb3\xf9\xe3\x9e\x45\x3b\x50\xa2\xd6\x0e\xa4\x33\x06\x2d\x45\x3b\x1c\xaa\x45\x91\x2e\xea\x73\x94\xa4\x53\x69\xf2\x54\xf6\xbd\xe1\x59\x40\x7f\xa9\x24\x26\x85\xe3\x93\x7b\xfb\xd1\x58\xff\x85\xd4\x0a\x54\x00\x01\xe4\x9c\x86\xab\x52\xc9\x12\xa4\xab\x75\xde\x07\x57\x70\xed\x6a\xb8\x52\x54\xf6\xc2\x83\xa9\x35\xa9\x4a\x23\x3c\x44\xbb\x50\x36\x6c\x3e\x6f\xc8\xf9\xfd\xf9\x67\xdb\xbc\x6d\x9b\xdb\xae\xb9\xeb\x3e\xfe\x68\x7f\x7d\x59\xdd\x7d\x5d\x7d\xfa\xd0\x7d\xbb\x69\x9b\xdb\xbf\x43\xb0\xba\x79\xd7\x2d\xbf\x77\xef\x97\x71\xf4\x27\x00\x00\xff\xff\x8f\x02\x7e\xb7\x05\x02\x00\x00")

func zh_cnLc_messagesK8sPoBytes() ([]byte, error) {
	return bindataRead(
		_zh_cnLc_messagesK8sPo,
		"zh_CN/LC_MESSAGES/k8s.po",
	)
}

func zh_cnLc_messagesK8sPo() (*asset, error) {
	bytes, err := zh_cnLc_messagesK8sPoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "zh_CN/LC_MESSAGES/k8s.po", size: 517, mode: os.FileMode(420), modTime: time.Unix(1574260489, 0)}
	a := &asset{bytes: bytes, info: info}
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
	if err != nil {
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
	"zh_CN/LC_MESSAGES/k8s.mo": zh_cnLc_messagesK8sMo,
	"zh_CN/LC_MESSAGES/k8s.po": zh_cnLc_messagesK8sPo,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"zh_CN": &bintree{nil, map[string]*bintree{
		"LC_MESSAGES": &bintree{nil, map[string]*bintree{
			"k8s.mo": &bintree{zh_cnLc_messagesK8sMo, map[string]*bintree{}},
			"k8s.po": &bintree{zh_cnLc_messagesK8sPo, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
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

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
