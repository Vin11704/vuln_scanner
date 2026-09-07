package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
	colorBlue = "\033[34m"
	colorRed = "\033[31m%s\033[0m"
	colorYellow = "\033[33m%s\033[0m"
	colorCyan = "\033[36m%s\033[0m"
)

func colorSeverity(sev string) string {
	switch strings.ToLower(sev) {
	case "critical", "high":
		return fmt.Sprintf(colorRed, sev) // Red
	case "medium":
		return fmt.Sprintf(colorYellow, sev) // Yellow
	case "low":
		return fmt.Sprintf(colorCyan, sev) // Cyan
	default:
		return sev
	}
}

func main() {
	fmt.Println("test Network Vulnerability Scanner")

	root := "../LLM_prompt_detector"
	scanDir := "../LLM_prompt_detector"

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
			fmt.Println(colorBlue + "FINAL VERDICT:\n " + colorReset + colorGreen + "SAFE " + colorReset + "-> All files are listed in .gitignore and will not be committed/leaked")
		}
		
	}
	
}
