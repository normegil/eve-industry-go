package model

type Blueprint struct {
	ItemID             int64 `json:"item_id"`
	MaterialEfficiency int   `json:"material_efficiency"`
	TimeEfficiency     int   `json:"time_efficiency"`
	Quantity           int   `json:"quantity"`
	Runs               int   `json:"runs"`
}

type APIBlueprint struct {
	TypeID             int64  `json:"type_id" bson:"type_id"`
	ItemID             int64  `json:"item_id" bson:"item_id"`
	LocationFlag       string `json:"location_flag" bson:"location_flag"`
	LocationID         int64  `json:"location_id" bson:"location_id"`
	MaterialEfficiency int    `json:"material_efficiency" bson:"material_efficiency"`
	TimeEfficiency     int    `json:"time_efficiency" bson:"time_efficiency"`
	Quantity           int    `json:"quantity" bson:"quantity"`
	Runs               int    `json:"runs" bson:"runs"`
}
