package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
)

func crawlPage(ctx context.Context, rawBaseURL, rawCurrentURL string, pages map[string]int) {
	select {
	case <-ctx.Done():
		fmt.Println("\nCrawling stopped by user")
		return
	default:
	}

	if !isSameDomain(rawBaseURL, rawCurrentURL) {
		return
	}
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, ok := pages[normalizedURL]
	if ok {
		pages[normalizedURL] += 1
		return
	} else {
		pages[normalizedURL] = 1
	}
	html, err := getHTML(rawBaseURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	urls, err := getURLsFromHTML(html, rawBaseURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, url := range urls {
		select {
		case <-ctx.Done():
			fmt.Println("\nCrawling stopped by user")
			return
		default:
			crawlPage(ctx, rawBaseURL, url, pages)
		}
	}
}

func isSameDomain(rawBaseURL, rawCurrentURL string) bool {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("failed to parse base URL: %w\n", err)
		return false
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to parse current URL: %w\n", err)
		return false
	}

	baseHost := baseURL.Hostname()
	currentHost := currentURL.Hostname()

	return baseHost == currentHost
}
