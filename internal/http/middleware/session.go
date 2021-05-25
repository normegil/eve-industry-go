package middleware

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/http"
	"github.com/normegil/evevulcan/internal/model"
	stdhttp "net/http"
	"time"
)

type SessionHandler struct {
	SessionManager *scs.SessionManager
	DB             *db.DB
	ErrHandler     http.ErrorHandler
	Handler        stdhttp.Handler
}

func (s SessionHandler) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var token string
	cookie, err := r.Cookie(s.SessionManager.Cookie.Name)
	if err == nil {
		token = cookie.Value
	}

	ctx, err := s.SessionManager.Load(r.Context(), token)
	if err != nil {
		s.ErrHandler.Handle(w, fmt.Errorf("could not load session: %w", err))
		return
	}
	r = r.WithContext(ctx)

	if err := s.handleAuthenticationAction(r); nil != err {
		s.ErrHandler.Handle(w, err)
		return
	}

	identityIDAsInterface := s.SessionManager.Get(ctx, http.KeySessionIdentity)
	if nil != identityIDAsInterface {
		identityID := identityIDAsInterface.(int64)
		if model.IdentityAnonymous().ID != identityID {
			identity, err := s.DB.LoadIdentity(identityID)
			if err != nil {
				s.ErrHandler.Handle(w, fmt.Errorf("could not load identity '%d': %w", identityID, err))
				return
			}
			ctx = context.WithValue(ctx, http.KeyIdentity, *identity)
		}
	}
	sr := r.WithContext(ctx)

	switch s.SessionManager.Status(sr.Context()) {
	case scs.Unmodified:
		fallthrough
	case scs.Modified:
		token, expiry, err := s.SessionManager.Commit(ctx)
		if err != nil {
			s.ErrHandler.Handle(w, fmt.Errorf("could not commit session: %w", err))
			return
		}
		s.writeSession(w, token, expiry)
	case scs.Destroyed:
		s.writeSession(w, "", time.Time{})
	}

	s.Handler.ServeHTTP(w, sr)
}

func (s SessionHandler) handleAuthenticationAction(r *stdhttp.Request) error {
	authenticationAction := r.Header.Get("X-Authentication-Action")
	if authenticationAction != "" {
		identitySessionUpdater := AuthenticatedIdentitySessionUpdater{SessionManager: s.SessionManager}
		switch authenticationAction {
		case "sign-out":
			err := identitySessionUpdater.SignOut(r)
			if err != nil {
				return fmt.Errorf("couldn't sign out: %w", err)
			}
		default:
			return fmt.Errorf("unrecognized authentication action: '%s'", authenticationAction)
		}
	}
	return nil
}

func (s SessionHandler) writeSession(w stdhttp.ResponseWriter, token string, expiry time.Time) {
	cookie := &stdhttp.Cookie{
		Name:     s.SessionManager.Cookie.Name,
		Value:    token,
		Path:     s.SessionManager.Cookie.Path,
		Domain:   s.SessionManager.Cookie.Domain,
		Secure:   s.SessionManager.Cookie.Secure,
		HttpOnly: s.SessionManager.Cookie.HttpOnly,
		SameSite: s.SessionManager.Cookie.SameSite,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	} else if s.SessionManager.Cookie.Persist {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)        // Round up to the nearest second.
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1) // Round up to the nearest second.
	}

	w.Header().Add("Set-Cookie", cookie.String())
	addHeaderIfMissing(w, "Cache-Control", `no-cache="Set-Cookie"`)
	addHeaderIfMissing(w, "Vary", "Cookie")
}

func addHeaderIfMissing(w stdhttp.ResponseWriter, key, value string) {
	for _, h := range w.Header()[key] {
		if h == value {
			return
		}
	}
	w.Header().Add(key, value)
}

type AuthenticatedIdentitySessionUpdater struct {
	SessionManager *scs.SessionManager
}

//func (a AuthenticatedIdentitySessionUpdater) RenewSessionOnAuthenticatedIdentity(r *stdhttp.Request, username string) error {
//	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
//		return fmt.Errorf("could not renew session token: %w", err)
//	}
//	a.SessionManager.Put(r.Context(), KeySessionIdentity, username)
//	return nil
//}

func (a AuthenticatedIdentitySessionUpdater) SignOut(r *stdhttp.Request) error {
	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
		return fmt.Errorf("could not renew session token: %w", err)
	}
	a.SessionManager.Put(r.Context(), http.KeySessionIdentity, model.IdentityAnonymous().ID)
	return nil
}
