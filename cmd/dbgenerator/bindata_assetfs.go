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

var _templateDbschemaGotpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x1b\x6b\x6f\xe3\xb8\xf1\xb3\xf5\x2b\xb8\xc6\x62\x41\x05\x5e\xe5\x3e\x14\xfd\x90\x83\x8b\x3a\xd9\xcb\x35\xc8\x5d\x36\x48\xf6\xb6\x05\x82\x20\xa0\xa5\x91\xa3\x5a\x12\x5d\x8a\x7e\x9d\xe1\xff\x5e\x90\x14\x25\xea\x2d\xe7\x85\x5c\xaf\xc0\x62\x63\x91\xf3\x9e\xe1\x0c\x39\xa2\x8e\x8f\xd1\xdf\x67\x10\x03\x23\x1c\x3c\xf4\x85\xa2\x98\x72\x04\x5e\xc0\x11\x7f\x0c\x12\xe4\x07\x21\x8c\xd0\xfa\x31\x70\x1f\x51\x90\x20\xb2\xe4\x34\x22\x3c\x70\x49\x18\x6e\x51\x8e\x37\xdd\x22\xfe\x08\x7a\x80\x32\xc7\xb2\x16\xc4\x9d\x93\x19\xa0\xdd\xce\xb9\x56\x3f\xaf\x48\x04\xfb\xbd\x65\x05\xd1\x82\x32\x8e\xb0\x35\x18\xfa\x11\x1f\x5a\x83\xdd\x8e\x91\x78\x06\xe8\xe3\x7c\x84\x3e\xae\xd0\xc9\x18\x39\x17\x12\x26\xd9\xef\xad\xc1\x70\xb7\xfb\xb8\xda\xef\x25\x1c\xc4\x9e\xa0\x30\x18\xce\x02\xfe\xb8\x9c\x3a\x2e\x8d\x8e\xd7\x30\xdd\x7c\xe6\x74\x71\xec\xd2\x68\xd8\x30\xe5\x4d\x9b\x67\x8e\xc3\x60\x7a\xec\x13\x97\x53\xb6\xed\x07\x75\xbc\x20\xb3\x20\x26\x3c\xa0\x71\x13\x02\xb8\x8f\xb4\x6d\xee\x78\x41\x18\x89\x86\x96\x6d\x59\x7c\xbb\x00\x74\x1b\x06\x2e\x3c\xec\x76\xce\x2d\x67\x4b\x97\x2b\x53\xa1\xbb\xfb\xa3\xd2\x90\x65\xf9\xcb\xd8\x45\x38\xa9\xc7\xb0\xd1\x8d\xb0\x24\xf6\x63\x24\xe0\x70\x84\x52\x91\x9d\x5f\xa9\x07\xa1\x8d\x80\x31\xca\x90\xfe\xbb\xb3\x06\x3e\x65\xe8\x61\x84\xa4\xd5\x95\x17\x12\x31\x3c\x08\x7c\x01\x23\x46\xfd\x18\xaf\xec\x1f\xe5\xd3\x87\x31\x8a\x83\x50\xce\x0f\x18\xf0\x25\x8b\xc5\xb0\x35\x18\xec\x2d\xf1\x2f\x1d\x8a\x83\xd0\xea\x29\xe8\x0d\x59\x1b\xb2\x96\x95\x7d\x4f\xe2\xfe\xcc\xe8\x72\x71\xba\xc5\x73\xd8\x9e\x07\x10\x7a\x28\xe1\x2c\x88\x67\x36\x8a\xc8\xe2\x4e\xfd\xbe\xaf\xba\x4b\xf0\x96\x52\xb5\x42\xed\xf6\x99\x62\x8c\xae\xcb\xaa\x79\x11\x59\xc8\x31\xba\x76\x26\xc9\xaf\x64\x81\x6d\x6b\x30\x58\xcd\x61\x2b\xd5\x8d\xb8\x73\xbb\x60\x41\xcc\xb1\x00\xbc\xd3\xf2\xdd\xdb\xca\x2a\x0f\x23\x24\xe1\xd8\x9d\xc0\xb8\xff\x11\x7d\xd8\xa6\x06\x51\x03\x68\x5c\x13\x65\x42\x20\x61\x21\x03\x88\x2c\x16\x10\x7b\x38\x1d\x90\x82\xda\xa6\x15\x59\x0f\x1b\x5e\xc2\xb6\xc3\x82\x7d\xec\xf7\x66\xd6\xcb\x95\x67\x74\x7d\xa8\xae\x93\xe4\xf2\x7b\x59\xd5\x11\x5a\x91\x70\x09\x45\xed\x65\x22\x70\x6e\x39\x65\x90\x69\x6b\x8c\xbd\x81\x76\x72\x2e\x97\xec\xfe\x50\x55\xbf\x31\x12\x27\x3e\x65\x11\xe6\xf2\x17\xb0\xc4\xf4\x97\x52\xe6\x5b\x3a\x65\xa3\xbb\xfb\x3a\x95\x23\x32\x07\x5c\x98\x1a\x85\x10\xe3\xc4\xb6\x95\xfe\x81\xb7\xa9\xb5\x00\xbb\x0b\xbc\x4d\xea\x23\x6d\x02\xa7\x46\xa2\x6a\xb0\x76\xa9\x75\xce\x68\xf4\x4b\x90\x70\xec\x11\x4e\x50\x10\x73\x60\x3e\x71\x61\xb7\xb7\x1b\xb2\xf5\xce\x1a\x48\x2b\x26\x23\x44\xe7\x42\x4e\x81\xe8\xe0\xea\xf2\xb2\x2d\xb1\x32\x3f\xd0\xb9\xd4\x40\xa7\x34\x81\x9a\x6b\xa7\x71\x0b\x7c\xd5\xc2\x55\x46\xf8\x54\xb7\x66\xc5\xa4\x23\x04\xbf\xa1\x6b\x2c\x29\x3a\xd8\x70\x85\x49\x4c\xc4\xc0\x20\xc9\x57\x76\xa2\xd7\xb4\x5a\xf8\xca\x50\x89\xb4\x5a\x11\x4a\xe9\xe8\x38\x8e\x6d\x0d\xac\x1c\x50\xc7\xc9\x15\xac\x4b\x92\x61\x97\x6f\x90\x28\x7a\xce\x19\x8d\x39\x6c\xb8\x5d\xc9\xf4\x42\xb3\xa8\x49\xab\xc8\xb9\x05\x9e\xa2\x0a\x5a\x76\xc6\x35\x42\x82\xed\xf1\x31\xaa\x90\xd3\xcf\x67\x34\x8a\x20\xe6\xfb\xbd\x2a\xb3\x65\xb8\x44\x3e\x08\xee\x53\x92\x00\x42\x28\xab\x97\xa7\x24\x01\x6b\x40\xa7\xff\x06\x97\x27\x75\x85\xb8\x6e\xd7\xa2\x00\x26\x9c\xb3\x60\xba\xe4\x90\x28\xb0\x8f\x2b\xf4\x59\xfd\x52\xdb\x17\x25\xf2\x67\x24\x59\x0a\x9b\x89\x8d\x84\x0e\x47\x52\x53\x05\x65\x38\x63\x3b\x93\x4d\x3e\x13\x89\x06\xb2\x2a\xa6\xe6\x20\x8e\x20\xe9\xa4\xe0\xf9\xca\xad\xa3\xf9\x5b\x02\x6a\x71\xd4\x53\xb5\x8b\x1b\x07\xc1\x24\xa5\x9e\x21\xe6\x6e\x20\xed\xac\x8a\xde\x2b\x45\x42\x13\x9b\x26\x97\x77\xf0\xfa\x69\x05\x31\xff\x7a\x85\x69\x8c\x1c\xc7\x99\x52\x1a\x36\xb3\xc8\x61\x55\x30\x1f\xc2\xe1\xfc\x1c\x53\xdf\xef\xcb\x43\x41\x1f\xc2\x45\xeb\x6e\x17\xcc\x55\xf5\x75\x06\xd7\xc7\x05\xf1\xc5\x17\xec\xca\x3f\x22\x9d\x75\xd9\x3e\x87\x3e\x44\x6a\x81\x65\x0b\xf2\xb5\xb2\xca\xd9\x4e\x51\xc5\x4f\x86\x63\xf1\x3f\x52\x80\xa5\x5d\xac\xae\x9d\x2d\x0a\x18\x34\x7a\x8b\xaf\x70\x6c\xb5\x1d\xad\x65\x59\x55\x2a\xc5\xe9\xd4\xe9\x5a\x14\x36\x2c\xcb\x1b\x3a\xd2\xb4\xe5\x60\xab\x1a\x06\x5a\x6f\x35\x14\x4e\xb4\x56\x7a\x78\x53\xe7\x06\x92\x65\xc8\x6d\x94\xfd\x1c\x21\xc2\x66\x89\x88\xde\x42\x79\x29\xca\x25\x04\x09\x7c\xad\xa7\x22\x6a\xa3\x71\xbe\xa7\xce\x0c\x71\x05\xeb\x74\x5a\x88\xfc\x6b\xe0\x79\x21\xac\x09\x03\x1c\xad\xe5\xc8\x84\xcd\x12\x2c\x38\xaa\x15\xb0\x2f\xdb\xf0\x00\x64\x9d\x3b\xdd\x25\x63\x10\xf3\x5e\xe9\xf3\x0a\xd6\x38\xc9\x46\xb2\xbd\x58\xba\x12\x94\x11\x6a\x9c\x50\xf2\x73\x81\x88\xc6\xd6\x32\x35\x33\xff\xaa\x2a\x08\xb6\x6b\x6a\x48\x66\x61\x5d\x66\xaa\xc6\x15\x87\x93\x82\xc1\x52\xd0\xbb\x93\xfb\x76\xbe\xff\xca\x19\x37\xee\x56\x52\xa2\xb5\xf3\x98\x38\x19\x85\x0e\x15\xaf\x60\x9d\x33\xd3\x66\x94\xc7\x3b\xb3\x38\x7d\xaa\x65\xb3\xdb\xb7\xd3\xbe\x88\x03\x9e\x13\x3f\xaa\xb7\xa1\x61\xc0\x86\xc3\x8c\x96\x21\x83\xec\xd4\x48\xc7\x7b\x75\x4d\xa4\xb4\xf4\x78\x06\xab\x07\xbe\x80\x4f\x96\x21\x3f\x57\x8f\x32\x86\x2f\x62\x0f\x36\xb8\x94\x06\xe5\x8c\xaa\xd5\xc5\xc2\x6d\xab\x04\x1c\x86\x20\x43\x1b\x13\x99\x64\x1e\xd2\x09\x19\x9f\x98\x74\x65\x9c\x47\xca\xf8\x03\xae\xc9\x5b\xc3\xdd\xce\xf9\x46\xa6\x61\xda\x8e\x19\x76\xd0\x91\x4f\x8d\x84\x4c\xe0\x0e\x4a\xa9\x0e\x06\x9d\x3c\xb9\xe8\xc4\xfb\xa1\x12\xff\xff\x0c\xf8\xe3\x35\x03\x3f\xc8\xcc\x97\xc2\x62\x62\x17\x52\x89\x01\x98\x6d\x68\xb4\x96\xec\x67\xe0\x98\x38\xda\x26\x0a\xb7\xbd\x90\x5d\x4f\xc4\x2e\x1a\x27\x74\xc9\x5c\x28\x37\x51\x6a\xf2\xb5\xb1\x5d\x51\x38\x79\x65\xb6\xf3\x79\xe1\xf7\x7c\x5a\x45\x41\x3a\xab\xd4\x4a\x27\x53\x1d\x7b\x67\x7c\xa1\xdf\x13\xf2\x3d\x06\xc6\x54\x6b\xc5\xce\xb6\xc0\x27\xe3\xd4\x29\xea\x8c\x62\x6e\x63\x84\x39\xb1\x3a\x84\x08\x44\x01\xa8\x4b\x8d\x62\x20\xb2\xa1\xd0\xe5\x06\xdc\x15\x26\xb6\xf3\x35\x06\x79\x1e\x55\x44\xd0\x18\x29\xba\xa9\x52\xd2\x7b\xff\x59\x02\xdb\xaa\xa5\x75\xd2\x49\xd0\xd2\x6d\x9e\x31\xfa\x72\x7a\xe1\x9c\x07\x0c\x6e\x80\x78\x41\x3c\xc3\x64\x84\x72\x52\xd5\xf6\x8f\xc1\x51\xe1\xe7\xc0\x5a\xca\x92\x90\x9a\x53\x4e\xa3\xca\x18\xbc\x12\x5f\x23\x20\xdb\x1d\x26\xcf\x95\x0c\xdc\x95\x79\xae\x1c\xa1\x76\x1f\x2e\xc8\x0c\x46\x28\x09\x7e\x07\x81\xd5\xe4\x53\x49\x40\xee\xc0\xfe\xfa\x97\x91\xe1\xdd\xc0\x47\x92\x61\xa1\xc6\x88\x67\x44\x9c\x42\x96\x95\x4a\xb4\xf8\x3e\x2b\x45\xf5\xbe\xba\x26\x33\xc0\x42\x54\xf9\x74\x1b\xfc\x0e\x58\x88\x9c\xfb\x51\x70\xb5\x1d\x69\x01\xfb\x80\x18\xe8\x47\xd7\x32\x3a\x81\x4f\x8a\x11\xf1\x38\x52\x6d\xc2\xbd\x35\x70\x85\x9d\x53\x72\x46\xc8\x68\xe1\xab\x41\x92\xac\x03\xee\x3e\xa6\xfd\x49\x70\x57\x0e\x16\xe7\x4e\x65\x38\x57\xc4\x57\x4d\x09\x3b\x11\x27\xf1\xee\xe8\x1a\x35\x14\xea\xa3\x95\x3c\xcb\x4b\xea\x2f\x4e\xdc\xa0\x5d\xac\xec\xbd\xe9\xae\xec\x72\xc3\x55\x1b\xb5\x23\xa7\xd5\x37\x5c\x47\x28\x88\x17\x4b\x7e\x43\xd7\x32\xf6\x6b\x9a\x2b\xdd\x1d\xd9\x15\x61\x88\x09\x02\xb5\x3a\x4b\xb7\x86\x10\xe3\x8c\x91\x8d\xfe\x86\x7e\x50\x31\x22\xb0\xc6\x0d\xb6\xca\xe0\xef\x7e\xb8\x17\x91\x8d\x20\x4c\xa0\x1b\xad\xb0\xd1\x32\x9b\x54\x74\x9d\x38\x65\x2b\x74\xd4\xad\xda\x0e\xeb\x61\x26\xfb\x83\x1b\xac\x68\x81\x0e\x73\xf5\x6c\xd2\xf6\x30\x60\xa9\xa5\xf9\x07\x32\x58\xc1\x06\xa6\xf2\x1d\xb6\x13\x39\xf0\x74\xfb\xd5\xf7\x13\x38\xbc\x94\x51\x89\xf6\xfe\x8b\x59\xaa\x9e\x12\xf7\x85\x0b\x5a\x7f\xda\xff\x2f\x6a\xff\x2b\x45\xed\x22\x4e\x80\x71\x6c\x23\xbc\x98\x17\xd7\x4b\x71\x2b\xbe\xdb\x39\xa7\xe0\x53\x06\x0a\x21\x4d\x19\xad\x9b\x71\x2d\x30\x1e\xba\x0c\x08\x0f\xe2\xd9\x70\x84\xc8\x48\x38\xd1\xce\x5f\x90\x56\x5f\x89\x66\x8a\x2c\xe6\x23\x54\xdc\xd7\x0b\x54\x19\x97\x10\x7b\x62\x3f\xaf\xa5\x97\xf2\x4d\x7c\x0e\xac\x20\x9e\x11\x35\x9f\x3e\x1d\x24\x2c\x78\xa6\xac\x3d\xf7\xd0\xbf\x2d\x3c\xc2\xe1\x05\xce\x3d\x99\xb1\x15\xc5\xfd\xb3\x32\x86\xb6\x55\x2a\x5d\x96\x80\xca\x7a\x2f\xc5\x7c\xe6\x24\x93\x4a\xdb\xe1\x25\x23\xd4\x97\x7d\x1b\xb1\x54\x93\x92\x4c\xda\x17\x26\xe5\x3e\x8e\xd8\x3c\xc9\x13\xc4\xf7\xc1\xe5\xe0\x19\x29\xfe\x6d\x3c\xb3\x79\x05\xd7\x68\x65\xca\xeb\xa8\x53\x92\xee\xf3\x6a\x1f\x27\x1d\xb0\x6a\x4e\x55\xb9\x4f\x3a\x7c\xe6\x4b\x20\x74\x77\xaf\xf7\x41\x6f\xbb\x9c\xb4\xac\x4a\x03\x91\x83\x95\x40\x59\x1b\x19\xbc\x80\x9f\xd1\x70\x19\xc5\x89\xf1\x0e\x59\x0b\x2b\x76\x58\x0a\x21\x7b\x89\x1c\x7b\xb0\x49\xa9\xe4\xaf\x5a\x53\x2d\x65\x76\xca\x09\xde\x49\xe0\x7b\x34\x46\x2e\x8d\x9c\xdb\x98\xcc\xe1\x8c\x24\xa0\x28\x36\xc4\x4e\xba\xea\x4a\x11\x64\x10\x7d\xa9\x95\xde\x66\x99\xfe\xc1\x54\x10\x57\x87\x54\xa3\xb4\x87\xc4\xd7\xe6\x65\x03\xec\xad\xb2\x44\x26\xfc\x9f\x34\xe2\x7a\x25\xb0\x56\x23\xbd\x87\xe0\x93\x91\xd7\x27\xf0\x8a\x27\xbe\xe2\x7e\xac\x4f\xa6\xcb\x42\xc9\xe0\x9b\x48\x63\xd5\xdf\xb9\x90\xb7\x3d\x04\xd0\x89\xe2\x38\xb2\x06\xfb\x03\x8b\xec\x2b\xea\xd6\xbe\xc8\x4a\xba\x6e\xde\x40\xd9\x5e\x19\x64\xbe\x4a\x80\x37\x88\xd0\xbb\x5e\xdd\x02\x57\x2c\xbd\x09\x7f\x81\x4d\x86\x94\xa9\xb4\x05\x94\x47\xc1\x23\x62\x0d\xa2\xec\x7e\x8e\x02\x53\xcd\x02\x33\xaf\xe8\x7c\xa2\x32\x88\x2b\x47\xf3\xdc\xa1\x14\x2e\xa5\x8e\xfc\x82\x4e\x61\x05\x29\xdc\xfe\xa9\xe3\x53\xf4\x2a\xd5\xaa\xc1\x2e\x07\xec\x4d\xab\xc9\xa2\x55\xd4\xfe\x4b\xe9\x35\xc2\xab\xb3\x5a\xbd\x66\xb8\x6d\xfe\x2c\xf1\x76\xd0\x66\xbb\x6c\xa2\xe7\x94\xa9\xf6\xc8\xd3\x1e\xeb\x13\x82\xdf\xe5\x4d\xba\xae\x08\x84\x6d\xa2\x00\xd1\x91\x37\x75\x2e\xb3\xc7\x7e\xd9\xed\xb9\xb1\x95\xb1\xeb\x91\xcf\x32\x58\x47\x5e\xc3\xac\x7b\x5f\xd7\x10\x04\x06\xaa\x50\x10\xdb\x2f\x9e\x7a\x6a\xf4\x78\x76\xfe\xe9\x90\xba\x2b\x04\x64\x33\xe5\x29\x87\xe6\xf6\xde\x51\xb5\x8f\xd3\x7c\xfe\x94\x22\xa4\x4d\xd6\xf4\x8e\x3d\xaa\xd9\x4e\xb7\x45\x51\xe1\xa6\x8a\x79\x77\xb4\xcf\x81\x5a\xee\x09\x9a\xf8\xe7\xcd\xa5\x67\xf2\xaf\x6b\x88\xed\x9f\xd5\xc5\x0a\x7c\xb4\x98\x9b\x2d\xd2\x43\x8e\xe9\x66\xdf\xbe\x4f\x3f\x4c\xf6\xe7\x50\xcf\xcc\xf2\x05\x42\x78\x5a\x57\xac\xf1\x1c\xaf\x48\x3e\xa3\x50\xa5\x32\x35\x36\x5c\x3c\x31\xff\xc2\xbd\x30\xcd\xf3\xa0\x06\x98\x14\xe4\xf0\x06\x98\xe2\xf5\xaa\x0d\xb0\x17\xf2\x41\x4b\xd7\xeb\xc9\x4e\xe8\x55\x88\x33\xf6\x87\xb4\xba\x9a\xdd\xd1\xeb\x34\x78\x46\x97\xf1\xd3\x12\x6c\xe5\x7d\x53\x9b\x65\x15\x9b\x8e\x10\xf9\x69\x13\x24\xbc\xab\xd4\xd7\x0b\x33\xa5\x34\xec\x2b\x4b\xca\xa6\x43\x98\x1b\x48\x80\xe3\xfa\x5b\xf8\xbb\xdd\x67\x54\xb9\xd7\x2e\x11\xe4\x6d\x76\xe2\xec\x76\x1f\x57\xce\xe5\x7e\x8f\xc6\x48\xfe\xfc\x5e\xbc\xdc\xde\xf3\xc6\x92\xfa\x52\x83\xc6\x61\xda\x27\x12\x6a\x1f\xf0\x45\x4c\xfa\xfe\x34\xc7\x97\x17\x53\xd5\x1b\xd4\x5a\x0d\xe4\x5d\xb0\x49\xec\x29\x29\x14\xce\x5e\x7d\xe1\x34\x54\x6a\xec\xf7\x43\xf9\x9d\x93\x63\x6a\x95\xab\x95\x7f\x3a\x32\xc8\xbf\xc8\x29\xb5\x75\x0c\x75\x8c\xd7\x66\x0a\x48\x26\xfb\xc3\x44\x53\x6f\xa8\x72\xf1\x4e\x50\x87\xb0\x86\xb4\xfb\x86\x8f\x78\xea\x7c\xa1\xf7\x6f\x8c\xae\x1b\x4e\x3a\xb6\xfe\xf4\x6e\x0e\xdb\xca\x97\x2a\x02\xcd\x50\x77\x0e\xdb\xe7\x2a\x7b\x29\x95\xd5\xca\xe5\xa1\xf6\x0f\xc7\xa5\xf1\x0a\x18\xbf\xe5\x44\xd4\x6b\x15\x14\x93\x24\x9d\xf4\x08\x27\xdf\xb6\x0b\xf9\x86\x4f\x88\x68\x17\x91\x7e\x8a\xbd\x06\x3b\x75\xdd\xdd\x16\xbb\xc7\xe2\x8e\x4b\x99\xa0\xbc\x54\x77\x56\x66\x04\xa1\xec\x1c\xb6\x95\x37\xa5\xf5\xf6\x95\xef\x1d\xa5\x7d\xe7\x23\xb4\x32\x3e\x6e\x54\x5f\x08\x0d\xd4\x15\x41\xac\x66\xe5\x87\x3b\x42\x76\x4f\x5d\x34\x95\xc8\xe2\x4c\x87\x25\xe8\x7c\x8e\xf4\x71\x4e\x8c\x17\x5e\xcb\x8b\x21\x89\x1e\xf8\x68\x9e\x7e\x18\x28\xa5\x4c\x97\xde\x8f\x28\xf5\x9d\xa0\x32\x46\x73\xc9\xc9\xd8\xac\xa8\x61\xe3\xb3\xb2\x39\x6c\x33\x71\xf4\xa2\x54\xb6\xcf\x2e\x34\x48\x19\xc6\xca\x62\x77\x3f\xdc\x67\xd0\xda\x52\x5a\xc5\x03\xc3\xa5\xba\x38\x9e\x15\x2f\xab\xe6\x60\x31\xa2\xa5\x47\xb8\x4c\x12\xb1\x94\xde\x51\x62\xbb\x7c\xdf\x89\xed\xb2\x90\xd8\x1a\x84\x7d\x62\x62\xfb\x25\x48\xd4\xb5\x41\x97\xc6\x9e\x3c\x3d\x9f\xd1\x68\x41\x97\xb1\x97\x8c\x50\x42\x19\xaf\x56\xdb\xec\x43\xe3\x87\xec\x8a\x44\xfe\xb1\xb7\x73\x05\x6b\x41\x13\x18\x56\x9b\xf3\xf4\xec\xc2\x50\x5d\x3d\x37\xb7\x64\xcc\xf9\xca\x3c\x60\xa7\x5b\x2c\xd9\x66\x47\x1f\x21\x98\x33\x89\x3d\x6c\xdb\xce\xb5\xe0\x33\xc3\xa4\x70\x73\xd8\xf8\x6c\xb9\x97\xaa\x93\xa4\xe6\x2e\xce\xab\xa8\x2f\xd8\xbc\x1f\xfd\x4f\x09\x77\x1f\xbf\x93\x30\x90\xe7\xf4\xb6\xae\x9d\xa1\xa4\xc8\x82\x12\xd2\x38\xc8\xa5\x03\x88\x38\x6a\x25\xdb\xd5\x93\x82\x7c\xd7\x51\x64\x98\x5f\x2e\x4f\x5b\x86\x1d\x3b\xb0\x0c\xb1\xa3\x57\x6f\x08\x5b\x15\xa1\x96\xbb\x6f\x5c\xe3\x12\x42\xfc\x37\x00\x00\xff\xff\x69\x04\x1f\xfb\xe0\x41\x00\x00")

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

	info := bindataFileInfo{name: "template/dbschema.gotpl", size: 16864, mode: os.FileMode(420), modTime: time.Unix(1688570966, 0)}
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
