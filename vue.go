package main

import (
	"net/http"
)

type VueFileSystem struct {
	http.FileSystem
}

func (fs *VueFileSystem) Open(name string) (http.File, error) {
	open, err := fs.FileSystem.Open(name)
	if err != nil {
		return fs.FileSystem.Open("index.html")
	}
	return open, err
}
