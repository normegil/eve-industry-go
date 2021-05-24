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

func (a API) LoginURL() (*url.URL, error) {
	responseType := "code"
	scopes := []string{
		"esi-characters.read_blueprints.v1",
	}
	scopeAsStr := strings.Join(scopes, " ")
	urlStr := fmt.Sprintf("https://%s/oauth/authorize?response_type=%s&redirect_uri=%s&client_id=%s&scope=%s", a.SSODomainName, responseType, a.RedirectURL.String(), a.Client.ID, scopeAsStr)
	return url.Parse(urlStr)
}

func (a API) RequestIdentity(code string) (*model.Identity, *model.StoredAccessToken, error) {
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

	returnedIdentity := &model.Identity{
		ID:           identity.CharacterID,
		Name:         identity.CharacterName,
		RefreshToken: tokens.RefreshToken,
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

func (a API) identity(tokens model.Tokens) (*identityResponse, error) {
	identityURL := fmt.Sprintf("https://%s/oauth/verify", a.SSODomainName)
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
