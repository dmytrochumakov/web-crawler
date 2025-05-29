package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	} else {
		fmt.Printf("starting crawl of: %s\n", args[0])
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, stopping crawler...")
		cancel() // Cancel the context
	}()

	baseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfg := config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 5),
		wg:                 &sync.WaitGroup{},
	}
	cfg.wg.Add(1)
	go func() {
		defer cfg.wg.Done()
		cfg.crawlPage(ctx, baseURL.String())
	}()
	cfg.wg.Wait()
	fmt.Printf("Crawling completed. Found %d pages\n", len(cfg.pages))
}
