package eveapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const cacheExpirationTypeInfo = 7 * 24 * time.Hour
const cacheExpirationGroupInfo = 7 * 24 * time.Hour
const cacheExpirationCategoryInfo = 7 * 24 * time.Hour

type Universe struct {
	API
}

type APIType struct {
	Capacity        float64          `json:"capacity" bson:"capacity"`
	Description     string           `json:"description" bson:"description"`
	GraphicID       int32            `json:"graphic_id" bson:"graphic_id"`
	GroupID         int32            `json:"group_id" bson:"group_id"`
	IconID          int32            `json:"icon_id" bson:"icon_id"`
	MarketGroupID   int32            `json:"market_group_id" bson:"market_group_id"`
	Mass            float64          `json:"mass" bson:"mass"`
	Name            string           `json:"name" bson:"name"`
	PackagedVolume  float64          `json:"packaged_volume" bson:"packaged_volume"`
	PortionSize     int32            `json:"portion_size" bson:"portion_size"`
	Published       bool             `json:"published" bson:"published"`
	Radius          float64          `json:"radius" bson:"radius"`
	TypeID          int32            `json:"type_id" bson:"type_id"`
	Volume          float64          `json:"volume" bson:"volume"`
	DogmaAttributes []DogmaAttribute `json:"dogma_attributes" bson:"dogma_attributes"`
	DogmaEffects    []DogmaEffects   `json:"dogma_effects" bson:"dogma_effects"`
}

type DogmaAttribute struct {
	AttributeID int32   `json:"attribute_id" bson:"attribute_id"`
	Value       float64 `json:"value" bson:"value"`
}

type DogmaEffects struct {
	EffectID  int32 `json:"effect_id" bson:"effect_id"`
	IsDefault bool  `json:"is_default" bson:"is_default"`
}

func (a Universe) TypeByID(typeID int32) (APIType, error) {
	typeIDStr := strconv.FormatInt(int64(typeID), 10)
	data, err := a.request(RequestOptions{
		QueryID:         "type-" + typeIDStr,
		Method:          "GET",
		ExtraURL:        "/universe/types/" + typeIDStr,
		Body:            strings.NewReader(""),
		CacheExpiration: time.Now().Add(cacheExpirationTypeInfo),
	})
	if err != nil {
		return APIType{}, fmt.Errorf("requesting type '%d': %w", typeID, err)
	}

	var typeInfo APIType
	if err = json.Unmarshal(data, &typeInfo); nil != err {
		return APIType{}, fmt.Errorf("could not unmarshall type data for '%d': %w", typeID, err)
	}
	return typeInfo, nil
}

type APIGroup struct {
	GroupID    int32   `json:"group_id" bson:"group_id"`
	CategoryID int32   `json:"category_id" bson:"category_id"`
	Name       string  `json:"name" bson:"name"`
	Published  bool    `json:"published" bson:"published"`
	Types      []int32 `json:"types" bson:"types"`
}

func (a Universe) GroupByID(id int32) (APIGroup, error) {
	idStr := strconv.FormatInt(int64(id), 10)
	data, err := a.request(RequestOptions{
		QueryID:         "group-" + idStr,
		Method:          "GET",
		ExtraURL:        "/universe/groups/" + idStr,
		Body:            strings.NewReader(""),
		CacheExpiration: time.Now().Add(cacheExpirationGroupInfo),
	})
	if err != nil {
		return APIGroup{}, fmt.Errorf("requesting type '%d': %w", id, err)
	}

	var groupInfo APIGroup
	if err = json.Unmarshal(data, &groupInfo); nil != err {
		return APIGroup{}, fmt.Errorf("could not unmarshall type data for '%d': %w", id, err)
	}
	return groupInfo, nil
}

type APICategory struct {
	CategoryID int32   `json:"category_id" bson:"category_id"`
	Name       string  `json:"name" bson:"name"`
	Published  bool    `json:"published" bson:"published"`
	Groups     []int32 `json:"groups" bson:"groups"`
}

func (a Universe) CategoryByID(id int32) (APICategory, error) {
	idStr := strconv.FormatInt(int64(id), 10)
	data, err := a.request(RequestOptions{
		QueryID:         "category-" + idStr,
		Method:          "GET",
		ExtraURL:        "/universe/categories/" + idStr,
		Body:            strings.NewReader(""),
		CacheExpiration: time.Now().Add(cacheExpirationCategoryInfo),
	})
	if err != nil {
		return APICategory{}, fmt.Errorf("requesting type '%d': %w", id, err)
	}

	var info APICategory
	if err = json.Unmarshal(data, &info); nil != err {
		return APICategory{}, fmt.Errorf("could not unmarshall type data for '%d': %w", id, err)
	}
	return info, nil
}
