package dao

import (
	"github.com/normegil/evevulcan/internal/dao/mappers"
	"github.com/normegil/evevulcan/internal/eveapi"
	"github.com/normegil/evevulcan/internal/model"
)

type Character struct {
	API eveapi.API
}

func (c Character) Blueprints(identity model.Identity) ([]model.Blueprint, error) {
	apiBlueprints, err := c.API.WithAuthentification(identity).Character().Blueprints()
	if err != nil {
		return nil, err
	}
	blueprints := make([]model.Blueprint, 0)
	for _, apiBlueprint := range apiBlueprints {
		type_, err := c.Type(apiBlueprint.TypeID)
		if err != nil {
			return nil, err
		}
		blueprints = append(blueprints, mappers.ToModelBlueprint(apiBlueprint, type_))
	}
	return blueprints, nil
}
