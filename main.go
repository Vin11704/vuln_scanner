package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func colorSeverity(sev string) string {
	switch strings.ToLower(sev) {
	case "critical", "high":
		return fmt.Sprintf("\033[31m%s\033[0m", sev) // Red
	case "medium":
		return fmt.Sprintf("\033[33m%s\033[0m", sev) // Yellow
	case "low":
		return fmt.Sprintf("\033[36m%s\033[0m", sev) // Cyan
	default:
		return sev
	}
}

func main() {
	fmt.Println("test Network Vulnerability Scanner")

	root := "../<file_name>"
	scanDir := "../<file_name>"

	if !CheckGitignore(root) {
		fmt.Println("WARNING: no .gitignore found in", root)
		fmt.Println("Secrets in ignored files could still leak.")
		os.Exit(1)
	}

	gi, err := NewGitignore(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read .gitignore: %v\n", err)
		os.Exit(1)
	}

	rules := []Rule{NewSecretsRule()}
	findings := scan(scanDir, rules, gi, 1)

	if len(findings) == 0{
		fmt.Println("No vulnerabilities found.")
	} else { // final output
		isAllSafe := true

		for _, f := range findings {
			fmt.Printf("Potential severity: %s\n%s:%d [%s] %s\n", colorSeverity(f.Severity), f.Filename, f.Line, f.RuleID, f.Message)
			if !gi.Match(f.Filename) {
				isAllSafe = false
			}
		}

		if isAllSafe{
			fmt.Println(colorGreen + "SAFE " + colorReset + "-> All files are listed in .gitignore and will not be committed/leaked")
		}
		
	}
	
}
