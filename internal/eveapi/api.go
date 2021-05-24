package eveapi

import (
	"github.com/normegil/evevulcan/internal/config"
	"net/url"
)

type API struct {
	SSODomainName string
	Client        config.ClientAuth
	RedirectURL   url.URL
}
