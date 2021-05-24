package http

import (
	"fmt"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/eveapi"
	"net/http"
	"net/url"
)

type authHandler struct {
	AppBaseURL   url.URL
	EveAPI       eveapi.API
	ErrorHandler errorHandler
	DB           *db.DB
}

func (a *authHandler) login(w http.ResponseWriter, r *http.Request) {
	loginURL, err := a.EveAPI.LoginURL()
	if err != nil {
		a.ErrorHandler.handle(w, err)
		return
	}
	http.Redirect(w, r, loginURL.String(), http.StatusFound)
}

func (a *authHandler) callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]

	identity, accessToken, err := a.EveAPI.RequestIdentity(code)
	if err != nil {
		a.ErrorHandler.handle(w, fmt.Errorf("requesting identity: %w", err))
		return
	}

	if err = a.DB.InsertOrUpdateIdentity(*identity); nil != err {
		a.ErrorHandler.handle(w, fmt.Errorf("inserting loaded identity: %w", err))
		return
	}
	if err = a.DB.ReplaceAccessToken(*accessToken); err != nil {
		a.ErrorHandler.handle(w, fmt.Errorf("replacing access token: %w", err))
		return
	}

	http.Redirect(w, r, a.AppBaseURL.String(), http.StatusFound)
}
