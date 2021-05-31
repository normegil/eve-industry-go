package http

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/eveapi"
	"net/http"
	"net/url"
)

type authHandler struct {
	AppBaseURL     url.URL
	EveAPI         eveapi.SSO
	ErrorHandler   ErrorHandler
	DB             *db.DB
	SessionManager *scs.SessionManager
}

func (a *authHandler) login(w http.ResponseWriter, r *http.Request) {
	loginURL, err := a.EveAPI.LoginURL()
	if err != nil {
		a.ErrorHandler.Handle(w, err)
		return
	}
	http.Redirect(w, r, loginURL.String(), http.StatusFound)
}

func (a *authHandler) callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]

	identity, accessToken, err := a.EveAPI.RequestIdentity(code)
	if err != nil {
		a.ErrorHandler.Handle(w, fmt.Errorf("requesting identity: %w", err))
		return
	}

	if err = a.DB.InsertOrUpdateIdentity(*identity); nil != err {
		a.ErrorHandler.Handle(w, fmt.Errorf("inserting loaded identity: %w", err))
		return
	}
	if err = a.DB.ReplaceAccessToken(*accessToken); err != nil {
		a.ErrorHandler.Handle(w, fmt.Errorf("replacing access token: %w", err))
		return
	}

	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
		a.ErrorHandler.Handle(w, fmt.Errorf("could not renew session token: %w", err))
		return
	}
	a.SessionManager.Put(r.Context(), KeySessionIdentityID, identity.ID)
	http.Redirect(w, r, a.AppBaseURL.String(), http.StatusFound)
}

func (a *authHandler) signout(w http.ResponseWriter, r *http.Request) {
	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
		a.ErrorHandler.Handle(w, fmt.Errorf("could not renew session token: %w", err))
		return
	}
	a.SessionManager.Put(r.Context(), KeySessionIdentityID, nil)

}
