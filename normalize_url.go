package main

import (
	"net/url"
)

func normalizeURL(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	return u.Host + u.Path, nil
}
