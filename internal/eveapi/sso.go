package eveapi

import (
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/model"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (a SSO) LoginURL() (*url.URL, error) {
	responseType := "code"
	scopes := []string{
		"esi-characters.read_blueprints.v1",
	}
	scopeAsStr := strings.Join(scopes, " ")
	urlStr := fmt.Sprintf("https://%s/oauth/authorize?response_type=%s&redirect_uri=%s&client_id=%s&scope=%s", a.DomainName, responseType, a.RedirectURL.String(), a.Client.ID, scopeAsStr)
	return url.Parse(urlStr)
}

func (a SSO) RequestIdentity(code string) (*model.Identity, *model.StoredAccessToken, error) {
	bodyCode, err := tokenRequestBodyByCode(code)
	if err != nil {
		return nil, nil, fmt.Errorf("creating token request body: %w", err)
	}

	tokens, err := a.tokenRequest(bodyCode)
	if err != nil {
		return nil, nil, fmt.Errorf("request tokens: %w", err)
	}

	identity, err := a.identity(*tokens)
	if err != nil {
		return nil, nil, fmt.Errorf("request identity: %w", err)
	}

	portraits, err := a.portraits(identity.CharacterID)
	if err != nil {
		return nil, nil, fmt.Errorf("request portraits: %w", err)
	}

	returnedIdentity := &model.Identity{
		ID:           identity.CharacterID,
		Name:         identity.CharacterName,
		RefreshToken: tokens.RefreshToken,
		Portraits:    portraits,
		Role:         "user",
	}
	accessToken := &model.StoredAccessToken{
		CharacterID: identity.CharacterID,
		AccessToken: tokens.AccessToken,
		ExpiresIn:   tokens.ExpiresIn,
	}
	return returnedIdentity, accessToken, nil
}

type identityResponse struct {
	CharacterID   int64  `json:"CharacterID"`
	CharacterName string `json:"CharacterName"`
}

type portraitsResponse struct {
	Url64  string `json:"px64x64"`
	Url128 string `json:"px128x128"`
	Url256 string `json:"px256x256"`
	Url512 string `json:"px512x512"`
}

func (a SSO) identity(tokens model.Tokens) (*identityResponse, error) {
	identityURL := fmt.Sprintf("https://%s/oauth/verify", a.DomainName)
	identityRequest, err := http.NewRequest("GET", identityURL, strings.NewReader(""))
	if err != nil {
		return nil, fmt.Errorf("creating identity url: %w", err)
	}
	identityRequest.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

	identityResp, err := http.DefaultClient.Do(identityRequest)
	if err != nil {
		return nil, fmt.Errorf("identity request: %w", err)
	}

	identityBody, err := io.ReadAll(identityResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read identity request body: %w", err)
	}

	identity := &identityResponse{}
	if err = json.Unmarshal(identityBody, identity); nil != err {
		return nil, fmt.Errorf("unmarshall identity: %w", err)
	}
	return identity, nil
}

func (a SSO) portraits(id int64) (model.Portraits, error) {
	portraitURL := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/portrait/", id)
	portraitRequest, err := http.NewRequest("GET", portraitURL, strings.NewReader(""))
	if err != nil {
		return model.Portraits{}, fmt.Errorf("creating portrait request for %d: %w", id, err)
	}

	resp, err := http.DefaultClient.Do(portraitRequest)
	if err != nil {
		return model.Portraits{}, fmt.Errorf("request portraits for %d: %w", id, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Portraits{}, fmt.Errorf("read portraits request body: %w", err)
	}

	portraitsResponse := &portraitsResponse{}
	if err = json.Unmarshal(body, portraitsResponse); nil != err {
		return model.Portraits{}, fmt.Errorf("unmarshall portraits: %w", err)
	}
	return model.Portraits{
		URL64:  portraitsResponse.Url64,
		URL128: portraitsResponse.Url128,
		URL256: portraitsResponse.Url256,
		URL512: portraitsResponse.Url512,
	}, nil
}
