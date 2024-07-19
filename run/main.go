package main

import (
	"fmt"
	"regexp"
	"strings"
)

func extractNimandTopik(message string) (string, string) {
	var nim, topik string
	// Handle non-breaking spaces
	message = strings.ReplaceAll(message, "\u00A0", " ")

	// Regex patterns to extract NIM and topik
	nimPattern := regexp.MustCompile(`(?i)nim\s+(\d+)`)
	topikPattern := regexp.MustCompile(`(?i)topik\s+(.+)`)

	// Find matches in the message
	nimMatch := nimPattern.FindStringSubmatch(message)
	topikMatch := topikPattern.FindStringSubmatch(message)

	// Extract NIM
	if len(nimMatch) > 1 {
		nim = nimMatch[1]
	}

	// Extract Topik
	if len(topikMatch) > 1 {
		topik = strings.TrimSpace(topikMatch[1])
		// Remove the word "poin" from topik if it exists
		topik = strings.ReplaceAll(topik, "poin", "")
		topik = strings.TrimSpace(topik)
	}

	fmt.Printf("Extracted NIM: %s, Topik: %s\n", nim, topik)
	return nim, topik
}

func main() {
	// Example message
	message := "NIM 123456789 topik Machine Learning poin"
	nim, topik := extractNimandTopik(message)
	fmt.Printf("Final Extracted NIM: %s, Topik: %s\n", nim, topik)
}
