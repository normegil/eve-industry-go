package config

import (
	"net/url"
	"os"
)

func FrontendBaseURL() *url.URL {
	baseURL, err := url.Parse(os.Getenv(appPrefix + "URL_FRONTEND"))
	if err != nil {
		panic(err)
	}
	return baseURL
}
