package model

import "time"

type Tokens struct {
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	AccessToken  string `json:"access_token" bson:"access_token"`
	ExpiresIn    int    `json:"expires_in" bson:"expires_in"`
}

type StoredAccessToken struct {
	CharacterID int64     `json:"character_id" bson:"character_id"`
	AccessToken string    `json:"access_token" bson:"access_token"`
	Created     time.Time `json:"created" bson:"created"`
}
