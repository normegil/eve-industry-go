package middleware

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/http"
	stdhttp "net/http"
)

type SessionIdentityHandler struct {
	SessionManager *scs.SessionManager
	DB             *db.DB
	ErrHandler     http.ErrorHandler
	Handler        stdhttp.Handler
}

func (s SessionIdentityHandler) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	identityIDInterface := s.SessionManager.Get(r.Context(), http.KeySessionIdentityID)
	ctx := r.Context()
	if nil != identityIDInterface {
		identityID := identityIDInterface.(int64)
		identity, err := s.DB.LoadIdentity(identityID)
		if err != nil {
			s.ErrHandler.Handle(w, fmt.Errorf("load identity from session '%d': %w", identityID, err))
			return
		}
		ctx = context.WithValue(r.Context(), http.KeyIdentity, identity)
	}
	s.Handler.ServeHTTP(w, r.WithContext(ctx))
}
