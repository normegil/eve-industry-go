package eveapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AuthentifiedCharacter struct {
	AuthentifiedAPI
}

const cacheExpirationBlueprints = 5 * time.Minute

type APIBlueprint struct {
	TypeID             int32  `json:"type_id" bson:"type_id"`
	ItemID             int64  `json:"item_id" bson:"item_id"`
	LocationFlag       string `json:"location_flag" bson:"location_flag"`
	LocationID         int64  `json:"location_id" bson:"location_id"`
	MaterialEfficiency int    `json:"material_efficiency" bson:"material_efficiency"`
	TimeEfficiency     int    `json:"time_efficiency" bson:"time_efficiency"`
	Quantity           int    `json:"quantity" bson:"quantity"`
	Runs               int    `json:"runs" bson:"runs"`
}

func (c AuthentifiedCharacter) Blueprints() ([]APIBlueprint, error) {
	data, err := c.request(RequestOptions{
		QueryID:         "character-blueprints-" + strconv.FormatInt(c.CharacterID, 10),
		Method:          "GET",
		ExtraURL:        "/characters/" + strconv.FormatInt(c.CharacterID, 10) + "/blueprints",
		Body:            strings.NewReader(""),
		CacheExpiration: time.Now().Add(cacheExpirationBlueprints),
	})
	if err != nil {
		return nil, fmt.Errorf("requesting blueprints for '%d': %w", c.CharacterID, err)
	}

	var blueprints []APIBlueprint
	if err = json.Unmarshal(data, &blueprints); nil != err {
		return nil, fmt.Errorf("unmarshall response into blueprints: %w", err)
	}
	return blueprints, nil
}
