package main

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	z := html.NewTokenizer(reader)
	var urls []string

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		_, hasAttr := z.TagName()
		_, val, _ := z.TagAttr()
		if hasAttr {
			if bytes.Contains(val, []byte("https://")) {
				urls = append(urls, string(val))
			} else {
				urls = append(urls, rawBaseURL+string(val))
			}
		}
	}
	return urls, nil
}
