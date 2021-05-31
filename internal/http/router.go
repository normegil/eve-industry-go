package http

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/eveapi"
	"net/http"
	"net/url"
)

func Routes(appBaseURL url.URL, frontend http.FileSystem, database *db.DB, api eveapi.SSO, sessionManager *scs.SessionManager) (http.Handler, error) {
	r := chi.NewRouter()

	errorHandler := ErrorHandler{}
	auth := &authHandler{
		AppBaseURL:     appBaseURL,
		EveAPI:         api,
		ErrorHandler:   errorHandler,
		DB:             database,
		SessionManager: sessionManager,
	}
	r.Get("/auth/login", auth.login)
	r.Get("/auth/callback", auth.callback)
	r.Get("/auth/sign-out", auth.signout)

	users := UsersHandler{ErrorHandler: errorHandler}
	r.Get("/api/users/current", users.current)

	r.Mount("/", http.FileServer(&vueFileSystem{FileSystem: frontend}))

	return r, nil
}
