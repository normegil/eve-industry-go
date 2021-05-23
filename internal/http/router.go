package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(frontend http.FileSystem) (http.Handler, error) {
	r := chi.NewRouter()

	r.Get("/auth/login", authLogin)
	r.Get("/auth/callback", authCallback)

	r.Mount("/", http.FileServer(&vueFileSystem{FileSystem: frontend}))

	return r, nil
}
