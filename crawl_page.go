package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
)

func (cfg *config) crawlPage(ctx context.Context, rawCurrentURL string) {
	rawBaseURL := cfg.baseURL.String()

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
	val, ok := cfg.getPage(normalizedURL)
	if ok {
		val += 1
		cfg.setPage(normalizedURL, val)
		return
	} else {
		cfg.setPage(normalizedURL, 1)
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
			cfg.wg.Add(1)
			go func(url string) {
				defer cfg.wg.Done()
				cfg.crawlPage(ctx, url)
			}(url)
		}
	}
}

func (cfg config) getPage(key string) (int, bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	val, ok := cfg.pages[key]
	return val, ok
}

func (cfg config) setPage(key string, value int) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[key] = value
}

func isSameDomain(rawBaseURL string, rawCurrentURL string) bool {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("failed to parse current URL: %w\n", err)
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
