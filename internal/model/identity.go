package model

type Identity struct {
	ID           int64  `json:"character_id"`
	Name         string `json:"name"`
	RefreshToken string `json:"refresh_token"`
}
