package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/normegil/evevulcan/internal/config"
	"net/http"
)

func Routes(frontend http.FileSystem) (http.Handler, error) {
	r := chi.NewRouter()

	redirectURL, err := config.EveSSORedirectURL()
	if err != nil {
		return nil, fmt.Errorf("loading eve sso rediret url: %w", err)
	}
	auth := authHandler{
		DomainName:   config.EveSSODomainName(),
		Client:       config.EveSSOClientAuth(),
		RedirectURL:  *redirectURL,
		ErrorHandler: errorHandler{},
	}
	r.Get("/auth/login", auth.login)
	r.Get("/auth/callback", auth.callback)

	r.Mount("/", http.FileServer(&vueFileSystem{FileSystem: frontend}))

	return r, nil
}
