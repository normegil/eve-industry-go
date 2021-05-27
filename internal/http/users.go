package http

import (
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/model"
	"net/http"
)

type UsersHandler struct {
	ErrorHandler ErrorHandler
}

func (u UsersHandler) current(w http.ResponseWriter, r *http.Request) {
	identityInterface := r.Context().Value(KeyIdentity)
	identity := identityInterface.(*model.Identity)
	if model.IdentityAnonymous().ID == identity.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	marshalledIdentity, err := json.Marshal(identity)
	if err != nil {
		u.ErrorHandler.Handle(w, fmt.Errorf("marshalling user: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(marshalledIdentity); err != nil {
		u.ErrorHandler.Handle(w, fmt.Errorf("writing user body: %w", err))
		return
	}
}
