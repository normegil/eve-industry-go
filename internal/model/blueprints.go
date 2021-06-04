package model

type Blueprint struct {
	ItemID             int64 `json:"item_id"`
	MaterialEfficiency int   `json:"material_efficiency"`
	TimeEfficiency     int   `json:"time_efficiency"`
	Quantity           int   `json:"quantity"`
	Runs               int   `json:"runs"`
	Type               Type  `json:"type"`
}
