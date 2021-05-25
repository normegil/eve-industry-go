package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/eveapi"
	"net/http"
	"net/url"
)

func Routes(appBaseURL url.URL, frontend http.FileSystem, database *db.DB, api eveapi.API) (http.Handler, error) {
	r := chi.NewRouter()

	auth := &authHandler{
		AppBaseURL:   appBaseURL,
		EveAPI:       api,
		ErrorHandler: errorHandler{},
		DB:           database,
	}
	r.Get("/auth/login", auth.login)
	r.Get("/auth/callback", auth.callback)

	r.Mount("/", http.FileServer(&vueFileSystem{FileSystem: frontend}))

	return r, nil
}
