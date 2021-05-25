package middleware

import (
	"context"
	"github.com/normegil/evevulcan/internal/http"
	"github.com/normegil/evevulcan/internal/model"
	stdhttp "net/http"
)

type AnonymousUserSetter struct {
	Handler stdhttp.Handler
}

func (a AnonymousUserSetter) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	r = r.WithContext(context.WithValue(r.Context(), http.KeyIdentity, model.IdentityAnonymous()))
	a.Handler.ServeHTTP(w, r)
}
