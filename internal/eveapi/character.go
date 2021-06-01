package eveapi

import (
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/model"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthentifiedCharacter struct {
	API
	CharacterID  int64
	RefreshToken string
}

func (c AuthentifiedCharacter) Blueprints() ([]model.APIBlueprint, error) {
	queryId := "character-blueprints-" + strconv.FormatInt(c.CharacterID, 10)
	data, err := c.DB.FromCache(queryId)
	if err != nil {
		return nil, fmt.Errorf("reading cache for '%s': %w", queryId, err)
	}
	if nil == data {
		accessToken, err := c.SSO.RequestAccessToken(c.CharacterID, c.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("retrieving access token: %w", err)
		}
		request, err := http.NewRequest("GET", c.BaseURL.String()+"/characters/"+strconv.FormatInt(c.CharacterID, 10)+"/blueprints", strings.NewReader(""))
		if err != nil {
			return nil, fmt.Errorf("create GET request to EVE api: %w", err)
		}
		request.Header.Add("Authorization", "Bearer "+accessToken)
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, fmt.Errorf("execute request: %w", err)
		}
		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("wrong response code '%d'", resp.StatusCode)
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}
		if err := c.DB.ToCache(db.APICacheObject{QueryID: queryId, Expiration: time.Time{}, Object: data}); nil != err {
			return nil, fmt.Errorf("saving response into cache: %w", err)
		}
	}

	var blueprints []model.APIBlueprint
	if err = json.Unmarshal(data, &blueprints); nil != err {
		return nil, fmt.Errorf("unmarshall response into blueprints: %w", err)
	}
	return blueprints, nil
}
