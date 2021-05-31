package eveapi

import (
	"github.com/normegil/evevulcan/internal/config"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/model"
	"net/url"
)

type SSO struct {
	DomainName  string
	Client      config.ClientAuth
	RedirectURL url.URL
}

type API struct {
	BaseURL url.URL
	SSO     SSO
	DB      *db.DB
}

func (a API) WithAuthentification(identity model.Identity) *AuthentifiedAPI {
	return &AuthentifiedAPI{
		API:          a,
		CharacterID:  identity.ID,
		RefreshToken: identity.RefreshToken,
	}
}

type AuthentifiedAPI struct {
	API
	CharacterID  int64
	RefreshToken string
}

func (a AuthentifiedAPI) Character() *AuthentifiedCharacter {
	return &AuthentifiedCharacter{
		CharacterID:  a.CharacterID,
		RefreshToken: a.RefreshToken,
		API:          a.API,
	}
}
