package eveapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/config"
	"io"
	"net/http"
)

type Tokens struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type tokenRequestBodyByCode struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}

func TokenRequestBodyByCode(code string) ([]byte, error) {
	body, err := json.Marshal(tokenRequestBodyByCode{
		GrantType: "authorization_code",
		Code:      code,
	})
	if err != nil {
		return nil, fmt.Errorf("code: marshalling token request body: %w", err)
	}

	return body, nil
}

type tokenRequestBodyByRefreshToken struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

func TokenRequestBodyByRefreshToken(refresh_token string) ([]byte, error) {
	body, err := json.Marshal(tokenRequestBodyByRefreshToken{
		GrantType:    "refresh_token",
		RefreshToken: refresh_token,
	})
	if err != nil {
		return nil, fmt.Errorf("refresh token: marshalling token request body: %w", err)
	}

	return body, nil
}

func TokenRequest(domainName string, client config.ClientAuth, body []byte) (*Tokens, error) {
	tokenURL := fmt.Sprintf("https://%s/oauth/token", domainName)
	tokenRequest, err := http.NewRequest("POST", tokenURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	tokenRequest.Header.Add("Content-Type", "application/json")
	tokenRequest.SetBasicAuth(client.ID, client.Secret)

	resp, err := http.DefaultClient.Do(tokenRequest)
	if err != nil {
		return nil, fmt.Errorf("request tokens: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read token response body: %w", err)
	}

	tokens := &Tokens{}
	if err = json.Unmarshal(bodyBytes, tokens); nil != err {
		return nil, fmt.Errorf("unmarshall token response: %w", err)
	}
	return tokens, nil
}
