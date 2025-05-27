package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	contentType := resp.Header.Get("Content-Type")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code is not 200, %d", resp.StatusCode)
	}

	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("incorrect content-type %s", contentType)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
