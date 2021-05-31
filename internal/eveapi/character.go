package eveapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type AuthentifiedCharacter struct {
	API
	CharacterID  int64
	RefreshToken string
}

func (c AuthentifiedCharacter) Blueprints() error {
	token, err := c.DB.AccessToken(c.CharacterID)
	if err != nil {
		return fmt.Errorf("retrieving access token: %w", err)
	}

	request, err := http.NewRequest("GET", c.BaseURL.String()+"/characters/"+strconv.FormatInt(c.CharacterID, 10)+"/blueprints", strings.NewReader(""))
	if err != nil {
		return fmt.Errorf("create GET request to EVE api: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}

}
