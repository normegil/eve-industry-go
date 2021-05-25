package model

type Identity struct {
	ID           int64  `json:"character_id" bson:"character_id"`
	Name         string `json:"name" bson:"name"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
