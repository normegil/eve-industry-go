package eveapi

import (
	"fmt"
	"github.com/normegil/evevulcan/internal/db"
	"github.com/normegil/evevulcan/internal/model"
	"io"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	BaseURL url.URL
	SSO     SSO
	DB      *db.DB
}

func (a API) Universe() *Universe {
	return &Universe{
		API: a,
	}
}

func (a API) WithAuthentification(identity model.Identity) *AuthentifiedAPI {
	return &AuthentifiedAPI{
		API:          a,
		CharacterID:  identity.ID,
		RefreshToken: identity.RefreshToken,
	}
}

type RequestOptions struct {
	QueryID         string
	Method          string
	ExtraURL        string
	Body            io.Reader
	CacheExpiration time.Time
}

func (a API) request(options RequestOptions) ([]byte, error) {
	data, err := a.DB.FromCache(options.QueryID)
	if err != nil {
		return nil, fmt.Errorf("reading cache for '%s': %w", options.QueryID, err)
	}
	if nil == data {
		request, err := http.NewRequest(options.Method, a.BaseURL.String()+options.ExtraURL, options.Body)
		if err != nil {
			return nil, fmt.Errorf("create %s request to '%s': %w", options.Method, options.ExtraURL, err)
		}
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
		if err := a.DB.ToCache(db.APICacheObject{QueryID: options.QueryID, Expiration: options.CacheExpiration, Object: data}); nil != err {
			return nil, fmt.Errorf("saving response into cache: %w", err)
		}
	}
	return data, nil
}

type AuthentifiedAPI struct {
	API
	CharacterID  int64
	RefreshToken string
}

func (a AuthentifiedAPI) Character() *AuthentifiedCharacter {
	return &AuthentifiedCharacter{
		AuthentifiedAPI: a,
	}
}

func (a AuthentifiedAPI) request(options RequestOptions) ([]byte, error) {
	data, err := a.DB.FromCache(options.QueryID)
	if err != nil {
		return nil, fmt.Errorf("reading cache for '%s': %w", options.QueryID, err)
	}
	if nil == data {
		accessToken, err := a.SSO.RequestAccessToken(a.CharacterID, a.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("create %s request to '%s': %w", options.Method, options.ExtraURL, err)
		}
		request, err := http.NewRequest(options.Method, a.BaseURL.String()+options.ExtraURL, options.Body)
		if err != nil {
			return nil, fmt.Errorf("create %s request to EVE api: %w", options.Method, err)
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
		if err := a.DB.ToCache(db.APICacheObject{QueryID: options.QueryID, Expiration: options.CacheExpiration, Object: data}); nil != err {
			return nil, fmt.Errorf("saving response into cache: %w", err)
		}
	}
	return data, nil
}
