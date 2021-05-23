package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var frontend embed.FS

func Frontend() (http.FileSystem, error) {
	distFS, err := fs.Sub(frontend, "dist")
	if err != nil {
		return nil, fmt.Errorf("load frontend subfolder '%s': %w", "dist", err)
	}
	return http.FS(distFS), nil
}
