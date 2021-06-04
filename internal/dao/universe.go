package dao

import (
	"github.com/normegil/evevulcan/internal/dao/mappers"
	"github.com/normegil/evevulcan/internal/model"
)

func (c Character) Type(id int32) (model.Type, error) {
	apiType, err := c.API.Universe().TypeByID(id)
	if err != nil {
		return model.Type{}, err
	}
	group, err := c.Group(apiType.GroupID)
	if err != nil {
		return model.Type{}, err
	}
	return mappers.ToModelType(apiType, group), nil
}

func (c Character) Group(id int32) (model.Group, error) {
	apiGroup, err := c.API.Universe().GroupByID(id)
	if err != nil {
		return model.Group{}, err
	}
	category, err := c.Category(apiGroup.CategoryID)
	if err != nil {
		return model.Group{}, err
	}
	return mappers.ToModelGroup(apiGroup, category), nil
}

func (c Character) Category(id int32) (model.Category, error) {
	apiCategory, err := c.API.Universe().CategoryByID(id)
	if err != nil {
		return model.Category{}, err
	}
	return mappers.ToModelCategory(apiCategory), nil
}
