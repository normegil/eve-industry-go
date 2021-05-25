package middleware

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

type RequestLogger struct {
	Handler http.Handler
}

func (l RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("method", r.Method).Str("url", r.RequestURI).Msg("request")
	l.Handler.ServeHTTP(w, r)
}
