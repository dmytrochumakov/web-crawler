package main

import (
	"fmt"
	"sort"
)

type Entry struct {
	link  string
	count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("\n=============================")
	fmt.Printf(" REPORT for %s\n", baseURL)
	fmt.Println("=============================")

	var entries []Entry
	for k, v := range pages {
		entries = append(entries, Entry{link: k, count: v})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count < entries[j].count
	})
	for _, entry := range entries {
		fmt.Printf("Found %d internal links to %s\n", entry.count, entry.link)
	}
}
