package model

type Tokens struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type StoredAccessToken struct {
	CharacterID int64  `json:"character_id"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
