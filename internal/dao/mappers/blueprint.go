package mappers

import (
	"github.com/normegil/evevulcan/internal/eveapi"
	"github.com/normegil/evevulcan/internal/model"
)

func ToModelBlueprint(apiType eveapi.APIBlueprint, type_ model.Type) model.Blueprint {
	return model.Blueprint{
		ItemID:             apiType.ItemID,
		MaterialEfficiency: apiType.MaterialEfficiency,
		TimeEfficiency:     apiType.TimeEfficiency,
		Quantity:           apiType.Quantity,
		Runs:               apiType.Runs,
		Type:               type_,
	}
}
