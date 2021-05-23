package http

import (
	"net/http"
)

type vueFileSystem struct {
	http.FileSystem
}

func (fs *vueFileSystem) Open(name string) (http.File, error) {
	open, err := fs.FileSystem.Open(name)
	if err != nil {
		return fs.FileSystem.Open("index.html")
	}
	return open, err
}
