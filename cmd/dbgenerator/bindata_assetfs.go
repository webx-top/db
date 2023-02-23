// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// template/dbschema.gotpl
// template/dbschema_init.gotpl
// template/model.gotpl
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
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


type assetFile struct {
	*bytes.Reader
	name            string
	childInfos      []os.FileInfo
	childInfoOffset int
}

type assetOperator struct{}

// Open implement http.FileSystem interface
func (f *assetOperator) Open(name string) (http.File, error) {
	var err error
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	content, err := Asset(name)
	if err == nil {
		return &assetFile{name: name, Reader: bytes.NewReader(content)}, nil
	}
	children, err := AssetDir(name)
	if err == nil {
		childInfos := make([]os.FileInfo, 0, len(children))
		for _, child := range children {
			childPath := filepath.Join(name, child)
			info, errInfo := AssetInfo(filepath.Join(name, child))
			if errInfo == nil {
				childInfos = append(childInfos, info)
			} else {
				childInfos = append(childInfos, newDirFileInfo(childPath))
			}
		}
		return &assetFile{name: name, childInfos: childInfos}, nil
	} else {
		// If the error is not found, return an error that will
		// result in a 404 error. Otherwise the server returns
		// a 500 error for files not found.
		if strings.Contains(err.Error(), "not found") {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
}

// Close no need do anything
func (f *assetFile) Close() error {
	return nil
}

// Readdir read dir's children file info
func (f *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	if len(f.childInfos) == 0 {
		return nil, os.ErrNotExist
	}
	if count <= 0 {
		return f.childInfos, nil
	}
	if f.childInfoOffset+count > len(f.childInfos) {
		count = len(f.childInfos) - f.childInfoOffset
	}
	offset := f.childInfoOffset
	f.childInfoOffset += count
	return f.childInfos[offset : offset+count], nil
}

// Stat read file info from asset item
func (f *assetFile) Stat() (os.FileInfo, error) {
	if len(f.childInfos) != 0 {
		return newDirFileInfo(f.name), nil
	}
	return AssetInfo(f.name)
}

// newDirFileInfo return default dir file info
func newDirFileInfo(name string) os.FileInfo {
	return &bindataFileInfo{
		name:    name,
		size:    0,
		mode:    os.FileMode(2147484068), // equal os.FileMode(0644)|os.ModeDir
		modTime: time.Time{}}
}

// AssetFile return a http.FileSystem instance that data backend by asset
func AssetFile() http.FileSystem {
	return &assetOperator{}
}

var _templateDbschemaGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x1b\x6b\x6f\xdb\x38\xf2\xb3\xf5\x2b\x58\xa3\x28\xa8\xc2\x55\xf6\xc3\xe1\x3e\x64\xe1\xc3\x39\xed\x76\xaf\xc8\x6e\x1a\x24\xdd\xde\x01\x41\x10\xd0\xd2\xc8\xd1\x59\x12\x7d\x14\xfd\x5a\x43\xff\xfd\x40\x52\x94\xa8\xb7\x9c\x17\xb2\xb7\x07\x14\x8d\x25\xcd\x7b\x86\x33\xc3\x11\x75\x72\x82\xfe\xbe\x80\x18\x18\xe1\xe0\xa1\x4f\x14\xc5\x94\x23\xf0\x02\x8e\xf8\x7d\x90\x20\x3f\x08\x61\x82\xb6\xf7\x81\x7b\x8f\x82\x04\x91\x35\xa7\x11\xe1\x81\x4b\xc2\x70\x8f\x0a\xbc\xf9\x1e\xf1\x7b\xd0\x37\x28\x73\x2c\x6b\x45\xdc\x25\x59\x00\x3a\x1c\x9c\x4b\xf5\xf3\x82\x44\x90\xa6\x96\x15\x44\x2b\xca\x38\xc2\xd6\x68\xec\x47\x7c\x6c\x8d\x0e\x07\x46\xe2\x05\xa0\xb7\xcb\x09\x7a\xbb\x41\xa7\x53\xe4\x7c\x91\x30\x49\x9a\x5a\xa3\xf1\xe1\xf0\x76\x93\xa6\x12\x0e\x62\x4f\x50\x18\x8d\x17\x01\xbf\x5f\xcf\x1d\x97\x46\x27\x5b\x98\xef\x3e\x70\xba\x3a\x71\x69\x34\x6e\x79\xe4\xcd\xdb\x9f\x9c\x84\xc1\xfc\xc4\x27\x2e\xa7\x6c\x3f\x0c\xea\x64\x45\x16\x41\x4c\x78\x40\xe3\x36\x04\x70\xef\x69\xd7\xb3\x93\x15\x61\x24\x1a\x5b\xb6\x65\xf1\xfd\x0a\xd0\x75\x18\xb8\x70\x77\x38\x38\xd7\x9c\xad\x5d\xae\x4c\x85\x6e\x6e\xdf\x57\x6e\x59\x96\xbf\x8e\x5d\x84\x93\x66\x0c\x1b\x5d\x09\x4b\x62\x3f\x46\x02\x0e\x47\x28\x13\xd9\xf9\x95\x7a\x10\xda\x08\x18\xa3\x0c\xe9\xbf\x07\x6b\xe4\x53\x86\xee\x26\x48\x5a\x5d\x79\x21\x11\xb7\x47\x81\x2f\x60\xc4\x5d\x3f\xc6\x1b\xfb\x47\x79\xf5\x66\x8a\xe2\x20\x94\xcf\x47\x0c\xf8\x9a\xc5\xe2\xb6\x35\x1a\xa5\x96\xf8\x97\xdd\x8a\x83\xd0\x1a\x28\xe8\x15\xd9\x1a\xb2\x56\x95\x7d\x4d\xe2\xfe\xcc\xe8\x7a\x75\xb6\xc7\x4b\xd8\x7f\x0e\x20\xf4\x50\xc2\x59\x10\x2f\x6c\x14\x91\xd5\x8d\xfa\x7d\x5b\x77\x97\xe0\x2d\xa5\xea\x84\x3a\xa4\xb9\x62\x8c\x6e\xab\xaa\x79\x11\x59\xc9\x7b\x74\xeb\xcc\x92\x5f\xc9\x0a\xdb\xd6\x68\xb4\x59\xc2\x5e\xaa\x1b\x71\xe7\x7a\xc5\x82\x98\x63\x01\x78\xa3\xe5\xbb\xb5\x95\x55\xee\x26\x48\xc2\xb1\x1b\x81\x71\xfb\x23\x7a\xb3\xcf\x0c\xa2\x6e\xa0\x69\x43\x94\x09\x81\x84\x85\x0c\x20\xb2\x5a\x41\xec\xe1\xec\x86\x14\xd4\x36\xad\xc8\x06\xd8\xf0\x1c\xf6\x3d\x16\x1c\x62\xbf\x17\xb3\x5e\xa1\x3c\xa3\xdb\x63\x75\x9d\x25\xe7\xdf\xab\xaa\x4e\xd0\x86\x84\x6b\x28\x6b\x2f\x13\x81\x73\xcd\x29\x83\x5c\x5b\xe3\xde\x0b\x68\x27\x9f\x15\x92\xdd\x1e\xab\xea\x37\x46\xe2\xc4\xa7\x2c\xc2\x5c\xfe\x02\x96\x98\xfe\x52\xca\x7c\xcb\x1e\xd9\xe8\xe6\xb6\x49\xe5\x88\x2c\x01\x97\x1e\x4d\x42\x88\x71\x62\xdb\x4a\xff\xc0\xdb\x35\x5a\x80\xdd\x04\xde\x2e\xf3\x91\x36\x81\xd3\x20\x51\x3d\x58\xfb\xd4\xfa\xcc\x68\xf4\x4b\x90\x70\xec\x11\x4e\x50\x10\x73\x60\x3e\x71\xe1\x90\xda\x2d\xd9\xfa\x60\x8d\xa4\x15\x93\x09\xa2\x4b\x21\xa7\x40\x74\x70\x7d\x79\xd9\x96\x58\x99\x6f\xe8\x52\x6a\xa0\x53\x9a\x40\x2d\xb4\xd3\xb8\x25\xbe\x6a\xe1\x2a\x23\xbc\x6b\x5a\xb3\xe2\xa1\x23\x04\xbf\xa2\x5b\x2c\x29\x3a\xd8\x70\x85\x49\x4c\xc4\xc0\x28\x29\x56\x76\xa2\xd7\xb4\x5a\xf8\xca\x50\x89\xb4\x5a\x19\x4a\xe9\xe8\x38\x8e\x6d\x8d\xac\x02\x50\xc7\xc9\x05\x6c\x2b\x92\x61\x97\xef\x90\x28\x7a\xce\x47\x1a\x73\xd8\x71\xbb\x96\xe9\x85\x66\x51\x9b\x56\x91\x73\x0d\x3c\x43\x15\xb4\xec\x9c\x6b\x84\x04\xdb\x93\x13\x54\x23\xa7\xaf\x3f\xd2\x28\x82\x98\xa7\xa9\x2a\xb3\x55\xb8\x44\x5e\x08\xee\x73\x92\x00\x42\x28\xaf\x97\x67\x24\x01\x6b\x44\xe7\xff\x06\x97\x27\x4d\x85\xb8\xa9\x6b\x51\x00\x33\xce\x59\x30\x5f\x73\x48\x14\xd8\xdb\x0d\xfa\xa0\x7e\xa9\xf6\x45\x89\xfc\x01\x49\x96\xc2\x66\xa2\x91\xd0\xe1\x48\x1a\xaa\xa0\x0c\x67\x6c\xe7\xb2\xc9\x6b\x22\xd1\x40\x56\xc5\xcc\x1c\xc4\x11\x24\x9d\x0c\xbc\x58\xb9\x4d\x34\x7f\x4b\x40\x2d\x8e\x66\xaa\x76\xb9\x71\x10\x4c\x32\xea\x39\x62\xe1\x06\xd2\xcd\xaa\xec\xbd\x4a\x24\xb4\xb1\x69\x73\x79\x0f\xaf\x9f\x36\x10\xf3\xaf\x17\x98\xc6\xc8\x71\x9c\x39\xa5\x61\x3b\x8b\x02\x56\x05\xf3\x31\x1c\x3e\x7f\xc6\xd4\xf7\x87\xf2\x50\xd0\xc7\x70\xd1\xba\xdb\x25\x73\xd5\x7d\x9d\xc3\x0d\x71\x41\xfc\xe5\x13\x76\xe5\x1f\x91\xce\xfa\x6c\x5f\x40\x1f\x23\xb5\xc0\xb2\x05\xf9\x46\x59\xe5\xd3\x5e\x51\xc5\x4f\x86\x63\xf1\x3f\x52\x80\x95\x2e\x56\xd7\xce\x0e\x05\x0c\x1a\x83\xc5\x57\x38\xb6\x6a\x47\x1b\x59\xd6\x95\xca\x70\x7a\x75\xba\x14\x85\x0d\xcb\xf2\x86\xde\x6b\xda\xf2\x66\xa7\x1a\x06\xda\x60\x35\x14\x4e\xb4\x55\x7a\x78\x73\xe7\x0a\x92\x75\xc8\x6d\x94\xff\x9c\x20\xc2\x16\x89\x88\xde\x52\x79\x29\xcb\x25\x04\x09\x7c\xad\xa7\x22\x6a\xa3\x69\xd1\x53\xe7\x86\xb8\x80\x6d\xf6\x58\x88\xfc\x6b\xe0\x79\x21\x6c\x09\x03\x1c\x6d\xe5\x9d\x19\x5b\x24\x58\x70\x54\x2b\x20\xad\xda\xf0\x08\x64\x9d\x3b\xdd\x35\x63\x10\xf3\x41\xe9\xf3\x02\xb6\x38\xc9\xef\xe4\xbd\x58\xb6\x12\x94\x11\x1a\x9c\x10\xf8\x48\x74\x20\xd9\x12\x40\x7f\x43\x3f\x98\x5a\x6b\xf0\x0b\xd8\x4a\x0c\x83\xc3\x44\xa1\xdc\xfc\x70\x6b\xcb\x6c\x59\x4e\xcb\x25\x03\x74\x51\xa9\xac\x9a\x66\x5a\x9d\x81\xf0\x55\x55\x2f\x6c\x37\xd4\xaf\xdc\xbb\xba\xc4\xd5\x1d\x2b\x36\x46\x25\x67\x65\xa0\x37\xa7\xb7\xdd\x7c\xff\x55\x30\x6e\xed\x94\x32\xa2\x8d\xcf\x31\x71\x72\x0a\x3d\x2a\x5e\xc0\xb6\x60\xa6\xad\x29\xb7\x96\x66\x61\x7c\xd7\xc8\xe6\x90\x76\xd3\xfe\x12\x07\xbc\x20\xfe\xbe\xd9\x86\x86\x01\x5b\x36\x52\x5a\x86\x1c\xb2\x57\x23\xbd\xd6\xea\xeb\xb1\x1e\x36\x0a\x56\xdf\xf8\x04\x3e\x59\x87\xfc\xb3\xba\x94\xeb\xe7\x4b\xec\xc1\x0e\xd7\x82\xe9\x1a\xb8\x0a\xa2\x4a\x44\xa9\xe4\x1f\x86\x20\x97\x15\x26\x32\xc1\xdd\x65\x0f\x54\x8c\x92\xbe\x6c\x77\x4f\x19\xbf\xc3\x0d\x39\x73\x7c\x38\x38\xdf\xc8\x3c\xcc\x46\x41\xe3\x1e\x3a\xf2\xaa\x95\x90\x09\xdc\x43\x29\xd3\xc1\xa0\x53\x24\x36\x9d\xf4\xdf\xd4\xe2\xff\x9f\x01\xbf\xbf\x64\xe0\x07\xb9\xf9\x32\x58\x4c\xca\xab\xd8\x00\xcc\x9b\x29\xad\x25\xfb\x19\x38\x26\x8e\xb6\x89\xc2\xed\x2e\xa2\x97\x33\xd1\xc1\xe3\x84\xae\x99\x0b\xd5\x01\x4e\x43\xad\x30\x5a\x25\x85\x53\x74\x05\x76\xf1\x5c\xf8\xbd\x78\xac\xa2\x20\x7b\xaa\xd4\xca\x1e\x66\x3a\x0e\xae\x36\x42\xbf\x07\xd4\x1a\x0c\x8c\xa9\xb1\x8e\x9d\xb7\xdf\xa7\xd3\xcc\x29\x6a\x7f\x64\xb6\x50\xc2\x9c\x58\x6d\x80\x04\xa2\x00\xd4\x65\x4e\x31\x10\xd5\x41\xe8\x72\x05\xee\x06\x13\xdb\xf9\x1a\x83\xdc\x0b\x2b\x22\x68\x8a\x14\xdd\x4c\x29\xe9\xbd\xff\xac\x81\xed\xd5\xd2\x3a\xed\x25\x68\xe9\x11\xd3\x14\x7d\x3a\xfb\xe2\x7c\x0e\x18\x5c\x01\xf1\x82\x78\x81\xc9\x04\x15\xa4\xea\xa3\x27\x83\xa3\xc2\x2f\x80\xb5\x94\x15\x21\x35\xa7\x82\x46\x9d\x31\x78\x15\xbe\x46\x40\x76\x3b\x4c\xee\x69\x19\xb8\x1b\x73\x4f\x3b\x41\xdd\x3e\x5c\x91\x05\x4c\x50\x12\xfc\x0e\x02\xab\xcd\xa7\x92\x80\xec\xfe\xfe\xfa\x97\x89\xe1\xdd\xc0\x47\x92\x61\xa9\xc6\x88\x6b\x44\x9c\x52\x96\x95\x4a\x74\xf8\x3e\x2f\x45\xcd\xbe\xba\x24\x0b\xc0\x42\x54\x79\x75\x1d\xfc\x0e\x58\x88\x5c\xf8\x51\x70\xb5\x1d\x69\x01\xfb\x88\x18\x18\x46\xd7\x32\xa6\x90\x0f\x8a\x11\x71\x39\x51\x23\xca\xd4\x1a\xb9\xc2\xce\x19\x39\x23\x64\xb4\xf0\xf5\x20\x49\xb6\x01\x77\xef\xb3\xd9\x28\xb8\x1b\x07\x8b\x3d\xaf\x32\x9c\x2b\xe2\xab\xa1\x84\x9d\x5a\xa3\x21\xd1\x35\x69\x29\xd4\xef\x37\x72\x8e\x20\xa9\x3f\x39\x71\x83\x76\xb9\xb2\x0f\xa6\xbb\xb1\xab\xc3\x5e\x6d\xd4\x9e\x9c\xd6\x3c\xec\x9d\xa0\x20\x5e\xad\xf9\x15\xdd\xca\xd8\x6f\x18\xec\xf4\x4f\x83\x37\x84\x21\x26\x08\x34\xea\x9c\x77\x9d\x39\x23\xa3\xf1\x14\x58\xd3\x16\x5b\xe5\xf0\xa2\xeb\xb4\x46\x29\x82\x30\x81\x7e\xb4\x52\xa3\x65\x0e\xc8\xe8\x36\x71\xaa\x56\xe8\xa9\x5b\x8d\xd3\xdd\xe3\x4c\xf6\x07\x37\x58\xd9\x02\x3d\xe6\x1a\x38\x20\x1e\x60\xc0\xca\x38\xf5\x0f\x64\xb0\x92\x0d\x4c\xe5\x7b\x6c\x27\x72\xe0\xd9\xfe\xab\xef\x27\x70\x7c\x29\xa3\x12\xed\xf5\x17\xb3\x4c\x3d\x25\xee\x13\x17\xb4\xe1\xb4\xff\x5f\xd4\xfe\x57\x8a\xda\x97\x38\x01\xc6\xb1\x8d\xf0\x6a\x59\x5e\x2f\xe5\x56\xfc\x70\x70\xce\xc0\xa7\x0c\x14\x42\x96\x32\x3a\x9b\x71\x2d\x30\x1e\xbb\x0c\x08\x0f\xe2\xc5\x78\x82\xc8\x44\x38\xd1\x2e\x5e\xce\xd6\x5f\xc7\xe6\x8a\xac\x96\x13\x54\xee\xeb\x05\xaa\x8c\x4b\x88\x3d\xd1\xcf\x6b\xe9\xa5\x7c\x33\x9f\x03\x2b\x89\x67\x44\xcd\xbb\x77\x47\x09\x0b\x9e\x29\xeb\xc0\x1e\xfa\xb7\x95\x47\x38\x3c\xc1\xbe\x27\x37\xb6\xa2\x98\x3e\x2a\x63\x68\x5b\x65\xd2\xe5\x09\xa8\xaa\xf7\x5a\x3c\xcf\x9d\x64\x52\xe9\xda\xbc\xe4\x84\x86\xb2\xef\x22\x96\x69\x52\x91\x49\xfb\xc2\xa4\x3c\xc4\x11\xbb\x07\x79\x82\xf8\x3e\xb8\x1c\x3c\x23\xc5\xbf\x8c\x67\x76\xcf\xe0\x1a\xad\x4c\x75\x1d\xf5\x4a\xd2\xbf\x5f\x1d\xe2\xa4\x23\x56\xcd\x99\x2a\xf7\x49\x8f\xcf\x7c\x09\x84\x6e\x6e\x75\x1f\xf4\xb2\xcb\x49\xcb\xaa\x34\x10\x39\x58\x09\x94\x8f\xb0\xc1\x0b\xf8\x47\x1a\xae\xa3\x38\x31\xde\x5f\x6b\x61\x45\x87\xa5\x10\xf2\x17\xd8\xb1\x07\xbb\x8c\x4a\xf1\x9a\x37\xd3\x52\x66\xa7\x82\xe0\x8d\x04\xbe\x45\x53\xe4\xd2\xc8\xb9\x8e\xc9\x12\x3e\x92\x04\x14\xc5\x96\xd8\xc9\x56\x5d\x25\x82\x0c\xa2\x4f\xb5\xd2\xbb\x2c\x33\x3c\x98\x4a\xe2\xea\x90\x6a\x95\xf6\x98\xf8\xda\x3d\x6d\x80\xbd\x54\x96\xc8\x85\xff\x93\x46\xdc\xa0\x04\xd6\x69\xa4\xd7\x10\x7c\x32\xf2\x86\x04\x5e\x79\xc7\x57\xee\xc7\x86\x64\xba\x3c\x94\x0c\xbe\x89\x34\x56\xf3\x79\x0f\x79\xd2\x44\x00\x9d\x2a\x8e\x13\x6b\x94\x1e\x59\x64\x9f\x51\xb7\xee\x45\x56\xd1\x75\xf7\x02\xca\x0e\xca\x20\xcb\x4d\x02\xbc\x45\x84\xc1\xf5\xea\x1a\xb8\x62\xe9\xcd\xf8\x13\x34\x19\x52\xa6\x4a\x0b\x28\xb7\x82\xef\x89\x35\x8a\xf2\xb3\x41\x0a\x4c\x0d\x0b\xcc\xbc\xa2\xf3\x89\xca\x20\xae\xbc\x5b\xe4\x0e\xa5\x70\x25\x75\x14\x87\x83\x4a\x2b\x48\xe1\x0e\x4f\x1d\xef\xa2\x67\xa9\x56\x2d\x76\x39\xa2\x37\xad\x27\x8b\x4e\x51\x87\x2f\xa5\xe7\x08\xaf\xde\x6a\xf5\x9c\xe1\xb6\xfb\xb3\xc4\xdb\x51\xcd\x76\xd5\x44\x8f\x29\x53\xdd\x91\xa7\x3d\x36\x24\x04\xbf\xcb\x53\x7c\x7d\x11\x08\xfb\x44\x01\xa2\xf7\xde\xdc\x39\xcf\x2f\x87\x65\xb7\xc7\xc6\x56\xce\x6e\x40\x3e\xcb\x61\x1d\x79\x04\xb4\xe9\x7d\x5d\x4b\x10\x18\xa8\x42\x41\x6c\x3f\x79\xea\x69\xd0\xe3\xd1\xf9\xa7\x47\xea\xbe\x10\x90\xc3\x94\x87\x6c\x9a\xbb\x67\x47\xf5\x39\x4e\xfb\xfe\x53\x8a\x90\x0d\x59\xb3\xf3\xfd\xa8\xa1\x9d\xee\x8a\xa2\xd2\x49\x15\xf3\xdc\xea\x90\x0d\xb5\xec\x09\xda\xf8\x17\xc3\xa5\x47\xf2\x6f\x1a\x88\xa5\x8f\x9a\x62\x05\x3e\x5a\x2d\xcd\x11\xe9\x31\xdb\x74\x73\x6e\x3f\x64\x1e\x26\xe7\x73\x68\x60\x66\xf9\x04\x21\x3c\x6c\x2a\xd6\xba\x8f\x57\x24\x1f\x51\xa8\x32\x99\x5a\x07\x2e\x9e\x78\xfe\xc4\xb3\x30\xcd\xf3\xa8\x01\x98\x14\xe4\xf8\x01\x98\xe2\xf5\xac\x03\xb0\x27\xf2\x41\xc7\xd4\xeb\xc1\x4e\x18\x54\x88\x73\xf6\xc7\x8c\xba\xda\xdd\x31\x68\x37\xf8\x91\xae\xe3\x87\x25\xd8\xda\xfb\xa6\x2e\xcb\x2a\x36\x3d\x21\xf2\xd3\x2e\x48\x78\x5f\xa9\x6f\x16\x66\x4e\x69\x38\x54\x96\x8c\x4d\x8f\x30\x57\x90\x00\xc7\xcd\x5f\x00\x1c\x0e\x1f\x50\xed\x4c\xbd\x44\x90\x27\xe9\x89\x73\x38\xbc\xdd\x38\xe7\x69\x8a\xa6\x48\xfe\xfc\x5e\x3e\x58\x3f\xf0\xc4\x92\xfa\x4a\x84\xc6\x61\x36\x27\x12\x6a\x1f\xf1\x35\x4e\xf6\xfe\xb4\xc0\x97\x87\x62\xd5\x1b\xd4\x46\x0d\xe4\x59\xb0\x59\xec\x29\x29\x14\x4e\xaa\xbe\xae\x1a\x2b\x35\xd2\x74\x2c\xbf\xb1\x72\x4c\xad\x0a\xb5\x8a\xcf\x56\x46\xc5\xd7\x40\x95\xb1\x8e\xa1\x8e\xf1\xda\x4c\x01\xc9\x64\x7f\x9c\x68\xea\x0d\x55\x21\xde\x29\xea\x11\xd6\x90\x36\x6d\xf9\x80\xa8\xc9\x17\xba\x7f\x63\x74\xdb\xb2\xd3\xb1\xf5\x67\x7f\x4b\xd8\xd7\xbe\x92\x11\x68\x86\xba\x4b\xd8\x3f\x56\xd9\x73\xa9\xac\x56\xae\x08\xb5\x7f\x38\x2e\x8d\x37\xc0\xf8\x35\x27\xa2\x5e\xab\xa0\x98\x25\xd9\x43\x8f\x70\xf2\x6d\xbf\x92\x6f\xf8\x84\x88\x76\x19\xe9\xa7\xd8\x6b\xb1\x53\xdf\xb9\x71\xd1\x3d\x96\x3b\x2e\x65\x82\xea\x52\x3d\x58\xb9\x11\x84\xb2\x4b\xd8\xd7\xde\x94\x36\xdb\x57\xbe\x77\x94\xf6\x5d\x4e\xd0\xc6\xf8\xb0\x52\x7d\x9d\x34\x52\x47\x04\xb1\x7a\x2a\x3f\x1a\x12\xb2\x7b\xea\xa0\xa9\x44\x16\x7b\x3a\x2c\x41\x97\x4b\xa4\xb7\x73\xe2\x7e\xe9\xb5\xbc\xb8\x25\xd1\x03\x1f\x2d\xb3\x8f\x12\xa5\x94\xd9\xd2\xfb\x11\x65\xbe\x13\x54\xa6\x68\x29\x39\x19\xcd\x8a\xba\x6d\x7c\xd2\xb6\x84\x7d\x2e\x8e\x5e\x94\xca\xf6\xf9\x81\x06\x29\xc3\x54\x59\xec\xe6\x87\xdb\x1c\x5a\x5b\x4a\xab\x78\x64\xb8\xd4\x17\xc7\xa3\xe2\x65\xd3\x1e\x2c\x46\xb4\x0c\x08\x97\x59\x22\x96\xd2\x2b\x4a\x6c\xe7\xaf\x3b\xb1\x9d\x97\x12\x5b\x8b\xb0\x0f\x4c\x6c\xbf\x04\x89\x3a\x36\xe8\xd2\xd8\x93\xbb\xe7\x8f\x34\x5a\xd1\x75\xec\x25\x13\x94\x50\xc6\xeb\xd5\x36\xff\xc8\xf9\x2e\x3f\x22\x51\x7c\x68\xee\x5c\xc0\x56\xd0\x04\x86\x55\x73\x9e\xed\x5d\x18\x6a\xaa\xe7\x66\x4b\xc6\x9c\xaf\xcc\x03\x76\xb6\xc7\x92\x6d\xbe\xf5\x11\x82\x39\xb3\xd8\xc3\xb6\xed\x5c\x0a\x3e\x0b\x4c\x4a\x27\x87\x8d\x4f\xa6\x07\xa9\x3a\x4b\x1a\xce\xe2\x3c\x8b\xfa\x82\xcd\xeb\xd1\xff\x8c\x70\xf7\xfe\x3b\x09\x03\xb9\x4f\xef\x9a\xda\x19\x4a\x8a\x2c\x28\x21\x8d\x8d\x5c\x76\x03\x11\x47\xad\x64\xbb\xbe\x53\x90\xef\x3a\xca\x0c\x8b\xc3\xe5\xd9\xc8\xb0\xa7\x03\xcb\x11\x7b\x66\xf5\x86\xb0\x75\x11\x1a\xb9\xfb\xc6\x31\x2e\x21\xc4\x7f\x03\x00\x00\xff\xff\xd2\xdb\x76\x5d\x5c\x42\x00\x00")

func templateDbschemaGotplBytes() ([]byte, error) {
	return bindataRead(
		_templateDbschemaGotpl,
		"template/dbschema.gotpl",
	)
}

func templateDbschemaGotpl() (*asset, error) {
	bytes, err := templateDbschemaGotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "template/dbschema.gotpl", size: 16988, mode: os.FileMode(420), modTime: time.Unix(1677166243, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templateDbschema_initGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8f\xcd\x8e\x9b\x30\x14\x85\xd7\xf1\x53\x1c\xb1\x02\x35\x85\x27\xa8\x54\x51\x36\xa8\x52\x14\x75\xd3\xb5\x81\x0b\x5c\x15\xec\xd4\x5c\x26\x41\x96\xdf\x7d\x64\x92\x49\x36\x33\x2b\xff\xe8\x3b\xf7\x3b\xb7\x28\xf0\x73\x20\x43\x4e\x0b\x75\xa8\x2c\x8c\x15\x50\xc7\x02\x8d\x9e\x27\x3a\xe2\x3a\x72\x3b\x82\x17\xe8\x55\xec\xac\x85\x5b\x3d\x4d\x1b\x5e\xa1\x66\x83\x8c\xf4\xf1\x61\x5d\xae\xd4\x45\xb7\xff\xf4\x40\xf0\x3e\x3f\xdf\xaf\x27\x3d\x53\x08\x4a\xf1\x7c\xb1\x4e\x90\xaa\x43\x32\xb0\x8c\x6b\x93\xb7\x76\x2e\xae\xd4\xdc\xbe\x8b\xbd\x14\x5d\x53\x4c\xdc\x14\xbd\x6e\xc5\xba\x2d\x51\x99\x52\x6f\xda\xe1\x2f\xcb\x78\x76\xd4\xf3\x0d\x3f\xd0\xaf\xa6\x4d\x45\x37\xd3\x3e\x14\x8b\x38\x36\x43\xf6\x38\xe1\xd5\xc1\x91\xac\xce\x20\x89\xf6\x3d\x14\x42\x82\x6f\x78\x46\x54\x50\xde\x73\x0f\xfa\x8f\xbc\x2a\x7f\xd3\x86\xbc\xa2\x5e\xaf\x93\xec\xaf\x10\x76\x67\x55\xd6\x51\x76\x6f\xf2\x02\x6a\xe5\x3d\x4d\x0b\x7d\x4a\x9d\xe8\x5a\x95\x75\x9a\x45\xc6\x74\x71\xdf\x58\x16\x6c\x58\xd2\xcc\x2b\x00\xd8\xcd\x86\xbe\x32\x47\xe4\xe9\x2c\xeb\x3f\x34\xf0\x22\xe4\xd2\xaa\xac\x8f\x71\xa1\x07\x97\x64\x8f\x61\x77\xcd\xc1\xfb\xbc\x36\x2c\xbf\x6c\x17\x8b\x05\xf5\x1e\x00\x00\xff\xff\x1e\xda\x5f\x0a\xd7\x01\x00\x00")

func templateDbschema_initGotplBytes() ([]byte, error) {
	return bindataRead(
		_templateDbschema_initGotpl,
		"template/dbschema_init.gotpl",
	)
}

func templateDbschema_initGotpl() (*asset, error) {
	bytes, err := templateDbschema_initGotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "template/dbschema_init.gotpl", size: 471, mode: os.FileMode(420), modTime: time.Unix(1613548584, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templateModelGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x93\xc1\x6f\xdb\x20\x14\xc6\xcf\xf0\x57\xbc\x45\x55\x05\x91\x47\xee\x99\x72\xe8\xaa\x1d\x26\x4d\x5d\xb4\x69\xa7\x69\x9a\x30\x10\x07\x25\x80\x85\x71\x93\x08\xf1\xbf\x4f\xe0\xa6\xf3\x9c\xa6\xd5\x6e\xf8\xf1\x78\xdf\xef\xfb\x0c\x2d\x17\x3b\xde\x28\x88\x91\xad\x87\xe5\x03\x37\x2a\x25\x8c\xb5\x69\x9d\x0f\x40\x30\x9a\xc5\xc8\xbe\x8b\xad\x32\xfc\xa9\x65\xcd\xc3\x36\xa5\x19\x46\xb3\x46\x87\x6d\x5f\x33\xe1\xcc\xe2\xa0\xea\xe3\xfb\xe0\xda\x85\x12\x5b\x77\x6d\x4f\xd6\xd7\x77\x16\x7b\x5d\x2f\x36\x5c\x04\xe7\x4f\x8b\x96\x37\xda\xf2\xa0\x9d\x9d\x61\x14\xa3\xe7\xb6\x51\x70\xb3\xab\xe0\xe6\x11\x96\x2b\x60\x9f\x0b\x5d\x97\x52\xc1\xbb\x79\x2c\x38\x31\x2a\x2b\x53\xc2\x14\x23\x84\x10\xde\xf4\x56\xc0\x83\x3a\x64\xfc\xe0\x7b\x11\x06\x6b\x44\x84\x23\x64\x48\x76\xef\x6c\x50\xc7\x40\x61\x3e\x69\x81\x88\x91\x57\xa1\xf7\x16\x6e\x27\x5b\x11\x23\x34\x29\x2d\x61\x1a\xd0\x50\x67\x2f\x6b\xd3\x0a\xa3\x84\x13\xc6\xe1\xd4\x96\xe0\xff\x51\xee\xca\x47\x06\x98\x5f\x99\x3a\x39\x91\x27\x15\xa7\xc4\x5c\xf8\xa0\x20\xb6\x4a\xec\x08\x05\xe5\xbd\xf3\x23\x5b\x56\xef\x5f\x3f\x78\x27\x25\xa1\x40\xda\x1d\x68\x1b\x94\xdf\x70\xa1\x62\xaa\xf2\x9c\x61\x16\xcd\xc3\xf4\xa6\x14\x96\x2b\x30\xec\x49\xe9\x43\xa9\xbc\x5b\x65\x85\xdc\x32\x12\x2c\xa7\xb3\xf7\x73\xcd\x4c\xbd\xb0\xa2\xfa\x3a\xd8\x27\xa9\x03\x31\x07\xc8\x1d\x44\xd6\xec\x9b\xea\xfa\x7d\xa0\xf0\xbc\xac\x80\xfb\xa6\x03\xc6\xd8\x88\x7c\x94\xc0\x7f\x40\xbf\xc9\xfb\xa3\x95\x3c\x28\x62\x0e\x83\x28\x63\xec\x0d\xfa\x2f\xba\x0b\x6b\xde\x28\x22\x9c\x95\x30\x97\x35\xbb\x77\xa6\x75\xbd\x95\x5d\x05\x5d\xbe\xd2\x17\xe0\xe4\xe7\xaf\xe9\x9c\x6a\xf4\x0f\xbc\x12\xe5\x49\x5c\x76\xc5\x84\xd1\xef\xea\x6c\xf6\xef\x8b\xca\xf7\x32\x73\x28\x4f\x2e\x1c\x55\x70\x9b\x07\x56\x43\xbc\x1e\x5e\x0a\x78\x9c\x90\x67\x5f\xbd\x54\xfe\xe3\x89\x14\xf8\x12\x00\x4a\x15\x64\x7b\xec\xce\x4a\x42\x29\x5b\x67\xe5\x86\x98\xf3\x7b\x23\x94\x3e\x47\x3a\x88\xe5\x9c\x13\xfe\x13\x00\x00\xff\xff\xea\xb4\x9d\x55\x8c\x04\x00\x00")

func templateModelGotplBytes() ([]byte, error) {
	return bindataRead(
		_templateModelGotpl,
		"template/model.gotpl",
	)
}

func templateModelGotpl() (*asset, error) {
	bytes, err := templateModelGotplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "template/model.gotpl", size: 1164, mode: os.FileMode(420), modTime: time.Unix(1643015743, 0)}
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
	"template/dbschema.gotpl":      templateDbschemaGotpl,
	"template/dbschema_init.gotpl": templateDbschema_initGotpl,
	"template/model.gotpl":         templateModelGotpl,
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
	"template": &bintree{nil, map[string]*bintree{
		"dbschema.gotpl":      &bintree{templateDbschemaGotpl, map[string]*bintree{}},
		"dbschema_init.gotpl": &bintree{templateDbschema_initGotpl, map[string]*bintree{}},
		"model.gotpl":         &bintree{templateModelGotpl, map[string]*bintree{}},
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
