package main

import (
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"log"
	"net/http"
)

//go:embed ui/web/dist/*
var frontend embed.FS

func main() {
	webFrontend, err := getWebFrontendFS()
	if err != nil {
		panic(fmt.Errorf("load frontend: %w", err))
	}

	r := chi.NewRouter()
	r.Mount("/", http.FileServer(&VueFileSystem{webFrontend}))

	server := http.Server{
		Addr:    ":18080",
		Handler: r,
	}
	log.Printf("server listening: %s", server.Addr)
	if err := server.ListenAndServe(); nil != err {
		panic(err)
	}
}

func getWebFrontendFS() (http.FileSystem, error) {
	uiFS, err := fs.Sub(frontend, "ui")
	if err != nil {
		return nil, fmt.Errorf("load frontend subfolder '%s': %w", "ui", err)
	}
	webFS, err := fs.Sub(uiFS, "web")
	if err != nil {
		return nil, fmt.Errorf("load frontend subfolder '%s': %w", "web", err)
	}
	distFS, err := fs.Sub(webFS, "dist")
	if err != nil {
		return nil, fmt.Errorf("load frontend subfolder '%s': %w", "dist", err)
	}
	return http.FS(distFS), nil
}
