package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/normegil/evevulcan/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Routes(frontend http.FileSystem, database *mongo.Database) (http.Handler, error) {
	r := chi.NewRouter()

	redirectURL, err := config.EveSSORedirectURL()
	if err != nil {
		return nil, fmt.Errorf("loading eve sso rediret url: %w", err)
	}
	auth := &authHandler{
		SSODomainName: config.EveSSODomainName(),
		Client:        config.EveSSOClientAuth(),
		RedirectURL:   *redirectURL,
		ErrorHandler:  errorHandler{},
		Mongo:         database,
	}
	r.Get("/auth/login", auth.login)
	r.Get("/auth/callback", auth.callback)

	r.Mount("/", http.FileServer(&vueFileSystem{FileSystem: frontend}))

	return r, nil
}
