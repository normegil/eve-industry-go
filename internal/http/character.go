package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/normegil/evevulcan/internal/dao"
	"github.com/normegil/evevulcan/internal/model"
	"net/http"
	"sort"
	"strconv"
)

type CharactersHandler struct {
	CharacterDAO dao.Character
	ErrorHandler ErrorHandler
}

func (c CharactersHandler) blueprints(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	sortParam := queryParams.Get("sort")
	if sortParam == "" {
		sortParam = "id"
	}
	offsetStr := queryParams.Get("page")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("wrong page param '%s': %w", offsetStr, err))
		return
	}
	offset -= 1 // Index start at 0
	limitStr := queryParams.Get("per_page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("wrong per_page param '%s': %w", limitStr, err))
		return
	}

	identityInterface := r.Context().Value(KeyIdentity)
	identity := identityInterface.(*model.Identity)
	if model.IdentityAnonymous().ID == identity.ID {
		c.ErrorHandler.Handle(w, Error{Code: 40100, Status: http.StatusUnauthorized, Err: fmt.Errorf("cannot load your blueprints: no characters identified")})
		return
	}

	blueprints, err := c.CharacterDAO.Blueprints(*identity)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("requesting blueprints for '%d': %w", identity.ID, err))
		return
	}

	sort.Slice(blueprints, func(i, j int) bool {
		return blueprints[i].ItemID < blueprints[j].ItemID
	})
	total := len(blueprints)
	startIndex := limit * offset
	if startIndex > total {
		startIndex = total
	}
	endIndex := limit * (1 + offset)
	if endIndex > total {
		endIndex = total
	}
	blueprintsToSend := blueprints[startIndex:endIndex]
	lastPage := total / limit
	if total%limit != 0 {
		lastPage += 1
	}

	resp := model.Collection{
		Total:     total,
		FromIndex: startIndex,
		ToIndex:   endIndex,
		PerPage:   limit,
		LastPage:  lastPage,
		Data:      blueprintsToSend,
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("marshalling blueprints: %w", err))
		return
	}

	if _, err = w.Write(respBody); nil != err {
		c.ErrorHandler.Handle(w, fmt.Errorf(": %w", err))
		return
	}
}

func (c CharactersHandler) ownedBlueprint(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemID")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("item ID could not be parsed '%s': %w", itemIDStr, err))
		return
	}

	identityInterface := r.Context().Value(KeyIdentity)
	identity := identityInterface.(*model.Identity)
	if model.IdentityAnonymous().ID == identity.ID {
		c.ErrorHandler.Handle(w, Error{Code: 40100, Status: http.StatusUnauthorized, Err: fmt.Errorf("cannot load your owned blueprint: no characters identified")})
		return
	}

	ownedBlueprint, err := c.CharacterDAO.OwnedBlueprint(itemID, *identity)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("requesting owned blueprint for '%d': %w", identity.ID, err))
		return
	}

	respBody, err := json.Marshal(ownedBlueprint)
	if err != nil {
		c.ErrorHandler.Handle(w, fmt.Errorf("marshalling owned blueprint: %w", err))
		return
	}

	if _, err = w.Write(respBody); nil != err {
		c.ErrorHandler.Handle(w, fmt.Errorf(": %w", err))
		return
	}
}
