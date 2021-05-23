package config

import (
	"net/url"
	"os"
)

const appPrefix = "EVE_VULCAN_"
const ssoPrefix = appPrefix + "EVE_SSO_"

func EveSSODomainName() string {
	return os.Getenv(ssoPrefix + "DOMAIN_NAME")
}

func EveSSOClientID() string {
	return os.Getenv(ssoPrefix + "CLIENT_ID")
}

func EveSSORedirectURL() (*url.URL, error) {
	return url.Parse(os.Getenv(ssoPrefix + "REDIRECT_URL"))
}
