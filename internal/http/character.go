package http

import (
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/eveapi"
	"github.com/normegil/evevulcan/internal/model"
	"net/http"
)

type CharactersHandler struct {
	API          eveapi.API
	ErrorHandler ErrorHandler
}

func (c CharactersHandler) blueprints(w http.ResponseWriter, r *http.Request) {
	identityInterface := r.Context().Value(KeyIdentity)
	identity := identityInterface.(*model.Identity)
	if model.IdentityAnonymous().ID == identity.ID {
		c.ErrorHandler.Handle(w, Error{Code: 40100, Status: http.StatusUnauthorized, Err: fmt.Errorf("cannot load your blueprints: no characters identified")})
		return
	}

	blueprints, err := c.API.WithAuthentification(*identity).Character().Blueprints()
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("requesting blueprints for '%d': %w", identity.ID, err))
		return
	}

	respBody, err := json.Marshal(blueprints)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("marshalling blueprints: %w", err))
		return
	}

	if _, err = w.Write(respBody); nil != err {
		c.ErrorHandler.Handle(w, fmt.Errorf(": %w", err))
		return
	}
}
