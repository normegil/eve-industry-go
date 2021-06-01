package eveapi

import (
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/config"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/model"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SSO struct {
	DomainName  string
	Client      config.ClientAuth
	RedirectURL url.URL
	DB          *db.DB
}

func (s SSO) LoginURL() (*url.URL, error) {
	responseType := "code"
	scopes := []string{
		"esi-characters.read_blueprints.v1",
	}
	scopeAsStr := strings.Join(scopes, " ")
	urlStr := fmt.Sprintf("https://%s/oauth/authorize?response_type=%s&redirect_uri=%s&client_id=%s&scope=%s", s.DomainName, responseType, s.RedirectURL.String(), s.Client.ID, scopeAsStr)
	return url.Parse(urlStr)
}

func (s SSO) RequestIdentity(code string) (*model.Identity, error) {
	bodyCode, err := tokenRequestBodyByCode(code)
	if err != nil {
		return nil, fmt.Errorf("creating token request body: %w", err)
	}

	tokens, err := s.tokenRequest(bodyCode)
	if err != nil {
		return nil, fmt.Errorf("request tokens: %w", err)
	}

	identity, err := s.identity(*tokens)
	if err != nil {
		return nil, fmt.Errorf("request identity: %w", err)
	}

	portraits, err := s.portraits(identity.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("request portraits: %w", err)
	}

	returnedIdentity := &model.Identity{
		ID:           identity.CharacterID,
		Name:         identity.CharacterName,
		RefreshToken: tokens.RefreshToken,
		Portraits:    portraits,
		Role:         "user",
	}
	if err = s.DB.InsertOrUpdateIdentity(*returnedIdentity); nil != err {
		return nil, fmt.Errorf("saving identity: %w", err)
	}

	accessToken := &model.StoredAccessToken{
		CharacterID: identity.CharacterID,
		AccessToken: tokens.AccessToken,
		ExpiresIn:   tokens.ExpiresIn,
	}
	if err = s.DB.ReplaceAccessToken(*accessToken); nil != err {
		return nil, fmt.Errorf("replacing access token: %w", err)
	}
	return returnedIdentity, nil
}

func (s SSO) RequestAccessToken(characterID int64, refreshToken string) (string, error) {
	token, err := s.DB.AccessToken(characterID)
	if err != nil {
		return "", fmt.Errorf("retrieving existing access token from cache: %w", err)
	}
	if token != nil {
		return token.AccessToken, nil
	}
	return s.requestNewAccessToken(characterID, refreshToken)
}

func (s SSO) requestNewAccessToken(characterID int64, refreshToken string) (string, error) {
	body, err := tokenRequestBodyByRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("constructing access token request body: %w", err)
	}
	tokens, err := s.tokenRequest(body)
	if err != nil {
		return "", fmt.Errorf("requesting new access token: %w", err)
	}
	if err = s.DB.ReplaceAccessToken(model.StoredAccessToken{CharacterID: characterID, AccessToken: tokens.AccessToken, ExpiresIn: tokens.ExpiresIn}); nil != err {
		return "", fmt.Errorf("could not store access token in cache: %w", err)
	}
	return tokens.AccessToken, nil
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

func (s SSO) identity(tokens model.Tokens) (*identityResponse, error) {
	identityURL := fmt.Sprintf("https://%s/oauth/verify", s.DomainName)
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

func (s SSO) portraits(id int64) (model.Portraits, error) {
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
