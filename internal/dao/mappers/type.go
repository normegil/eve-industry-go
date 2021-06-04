package mappers

import (
	"github.com/normegil/evevulcan/internal/eveapi"
	"github.com/normegil/evevulcan/internal/model"
)

func ToModelType(apiType eveapi.APIType, group model.Group) model.Type {
	return model.Type{
		ID:    apiType.TypeID,
		Name:  apiType.Name,
		Group: group,
	}
}

func ToModelGroup(apiType eveapi.APIGroup, category model.Category) model.Group {
	return model.Group{
		ID:       apiType.GroupID,
		Name:     apiType.Name,
		Category: category,
	}
}

func ToModelCategory(apiType eveapi.APICategory) model.Category {
	return model.Category{
		ID:   apiType.CategoryID,
		Name: apiType.Name,
	}
}
