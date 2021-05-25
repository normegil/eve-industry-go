package config

import (
	"net/url"
	"os"
)

func AppBaseURL() *url.URL {
	baseURL, err := url.Parse(os.Getenv(appPrefix + "URL"))
	if err != nil {
		panic(err)
	}
	return baseURL
}
