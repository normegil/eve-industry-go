package http

import (
	"encoding/json"
	"fmt"
	"github.com/normegil/evevulcan/internal/config"
	"github.com/normegil/evevulcan/internal/eveapi"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type authHandler struct {
	DomainName   string
	Client       config.ClientAuth
	RedirectURL  url.URL
	ErrorHandler errorHandler
}

func (a authHandler) login(w http.ResponseWriter, r *http.Request) {
	responseType := "code"
	scopes := []string{
		"esi-characters.read_blueprints.v1",
	}
	scopeAsStr := strings.Join(scopes, " ")
	loginURL := fmt.Sprintf("https://%s/oauth/authorize?response_type=%s&redirect_uri=%s&client_id=%s&scope=%s", a.DomainName, responseType, a.RedirectURL.String(), a.Client.ID, scopeAsStr)
	http.Redirect(w, r, loginURL, http.StatusFound)
}

type identityResponse struct {
	CharacterID   int64  `json:"CharacterID"`
	CharacterName string `json:"CharacterName"`
}

func (a authHandler) callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]

	bodyCode, err := eveapi.TokenRequestBodyByCode(code)
	if err != nil {
		a.ErrorHandler.handle(w, fmt.Errorf("creating token request body: %w", err))
		return
	}

	tokens, err := eveapi.TokenRequest(a.DomainName, a.Client, bodyCode)
	if err != nil {
		a.ErrorHandler.handle(w, fmt.Errorf("request tokens: %w", err))
		return
	}

	identity, err := requestIdentity(a.DomainName, *tokens)
	if err != nil {
		a.ErrorHandler.handle(w, fmt.Errorf("request identity: %w", err))
	}

	log.Printf("%+v", identity)
	// Store in MongoDB
}

func requestIdentity(domainName string, tokens eveapi.Tokens) (*identityResponse, error) {
	identityURL := fmt.Sprintf("https://%s/oauth/verify", domainName)
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
