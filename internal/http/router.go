package http

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/normegil/evevulcan/internal/dao"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/eveapi"
	"net/http"
	"net/url"
)

func Routes(appBaseURL url.URL, frontend http.FileSystem, database *db.DB, sso eveapi.SSO, sessionManager *scs.SessionManager) (http.Handler, error) {
	r := chi.NewRouter()

	errorHandler := ErrorHandler{}
	baseAPIURL, err := url.Parse("https://esi.evetech.net/latest")
	if err != nil {
		return nil, fmt.Errorf("parsing url base api url: %w", err)
	}
	api := eveapi.API{
		BaseURL: *baseAPIURL,
		SSO:     sso,
		DB:      database,
	}
	daos := dao.DAOs{API: api}

	auth := &authHandler{
		AppBaseURL:     appBaseURL,
		EveSSO:         sso,
		ErrorHandler:   errorHandler,
		DB:             database,
		SessionManager: sessionManager,
	}
	r.Get("/auth/login", auth.login)
	r.Get("/auth/callback", auth.callback)
	r.Get("/auth/sign-out", auth.signout)

	users := UsersHandler{ErrorHandler: errorHandler}
	r.Get("/api/users/current", users.current)

	characters := CharactersHandler{ErrorHandler: errorHandler, CharacterDAO: daos.Character()}
	r.Get("/api/characters/blueprints", characters.blueprints)

	r.Mount("/", http.FileServer(&vueFileSystem{FileSystem: frontend}))

	return r, nil
}
