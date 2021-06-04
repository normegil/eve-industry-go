package middleware

import (
	stdhttp "net/http"
	"net/url"
	"strings"
)

type CrossOriginRessourceSharing struct {
	Handler  stdhttp.Handler
	Frontend url.URL
}

func (c CrossOriginRessourceSharing) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	frontendURL := c.Frontend.String()
	frontendURL = strings.TrimSuffix(frontendURL, "/")
	w.Header().Add("Access-Control-Allow-Origin", frontendURL)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	c.Handler.ServeHTTP(w, r)
}
