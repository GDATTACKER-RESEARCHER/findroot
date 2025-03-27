package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// loadTLDs loads TLDs and second-level TLDs from a file.
func loadTLDs(tldFile string) (map[string]int, error) {
	tlds := make(map[string]int)
	file, err := os.Open(tldFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open TLD file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tld := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if tld != "" {
			tlds[tld] = strings.Count(tld, ".")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading TLD file: %w", err)
	}

	return tlds, nil
}

// extractRootDomains processes the domain list from stdin and extracts root domains.
func extractRootDomains(tlds map[string]int) (map[string]struct{}, error) {
	rootDomains := make(map[string]struct{})
	domainPattern := regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9.-]+`)

	scanner := bufio.NewScanner(os.Stdin)
	processedLines := 0

	for scanner.Scan() {
		processedLines++
		line := strings.ToLower(scanner.Text())
		domains := domainPattern.FindAllString(line, -1)

		for _, domain := range domains {
			bestMatch := ""
			bestLevel := -1
			for tld, level := range tlds {
				if strings.HasSuffix(domain, tld) {
					if level > bestLevel || (level == bestLevel && len(tld) > len(bestMatch)) {
						bestMatch = tld
						bestLevel = level
					}
				}
			}

			if bestMatch != "" {
				domainParts := strings.Split(domain, ".")
				if len(domainParts) > bestLevel {
					rootDomain := strings.Join(domainParts[len(domainParts)-bestLevel-1:], ".")
					rootDomains[rootDomain] = struct{}{}
				}
			}
		}

		if processedLines%10 == 0 {
			fmt.Printf("Processed %d lines. Root domains found: %d\r", processedLines, len(rootDomains))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	fmt.Printf("\nProcessing complete. Found %d root domains.\n", len(rootDomains))
	return rootDomains, nil
}

func main() {
	tldFile := "tld.txt"
	outputFile := flag.String("o", "root.txt", "File to save the extracted root domains")
	flag.Parse()

	fmt.Printf("Loading TLDs from %s...\n", tldFile)
	tlds, err := loadTLDs(tldFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading TLDs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d TLDs.\n", len(tlds))

	fmt.Println("Extracting root domains from input...")
	rootDomains, err := extractRootDomains(tlds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting root domains: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Saving root domains to %s...\n", *outputFile)
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for domain := range rootDomains {
		fmt.Fprintln(writer, domain)
	}
	writer.Flush()

	fmt.Printf("Root domains successfully saved to %s.\n", *outputFile)
}
