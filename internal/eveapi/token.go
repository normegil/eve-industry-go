package eveapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/model"
	"io"
	"net/http"
)

type tokenRequestBodyCode struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}

func tokenRequestBodyByCode(code string) ([]byte, error) {
	body, err := json.Marshal(tokenRequestBodyCode{
		GrantType: "authorization_code",
		Code:      code,
	})
	if err != nil {
		return nil, fmt.Errorf("code: marshalling token request body: %w", err)
	}

	return body, nil
}

type tokenRequestBodyRefreshToken struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

func tokenRequestBodyByRefreshToken(refresh_token string) ([]byte, error) {
	body, err := json.Marshal(tokenRequestBodyRefreshToken{
		GrantType:    "refresh_token",
		RefreshToken: refresh_token,
	})
	if err != nil {
		return nil, fmt.Errorf("refresh token: marshalling token request body: %w", err)
	}
	return body, nil
}

func (a API) tokenRequest(body []byte) (*model.Tokens, error) {
	tokenURL := fmt.Sprintf("https://%s/oauth/token", a.SSODomainName)
	tokenRequest, err := http.NewRequest("POST", tokenURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	tokenRequest.Header.Add("Content-Type", "application/json")
	tokenRequest.SetBasicAuth(a.Client.ID, a.Client.Secret)

	resp, err := http.DefaultClient.Do(tokenRequest)
	if err != nil {
		return nil, fmt.Errorf("request tokens: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read token response body: %w", err)
	}

	tokens := &model.Tokens{}
	if err = json.Unmarshal(bodyBytes, tokens); nil != err {
		return nil, fmt.Errorf("unmarshall token response: %w", err)
	}
	return tokens, nil
}
