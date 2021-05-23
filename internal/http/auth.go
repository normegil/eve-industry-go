package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type authHandler struct {
	DomainName  string
	ClientID    string
	RedirectURL url.URL
}

func (a authHandler) login(w http.ResponseWriter, r *http.Request) {
	responseType := "code"
	scopes := []string{
		"esi-characters.read_blueprints.v1",
	}
	scopeAsStr := strings.Join(scopes, " ")
	loginURL := fmt.Sprintf("https://%s/oauth/authorize?response_type=%s&redirect_uri=%s&client_id=%s&scope=%s", a.DomainName, responseType, a.RedirectURL.String(), a.ClientID, scopeAsStr)
	http.Redirect(w, r, loginURL, http.StatusFound)
}

func authCallback(w http.ResponseWriter, r *http.Request) {}
