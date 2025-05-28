package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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

	pages := make(map[string]int)
	crawlPage(ctx, args[0], args[0], pages)
	fmt.Printf("Crawling completed. Found %d pages\n", len(pages))
}
