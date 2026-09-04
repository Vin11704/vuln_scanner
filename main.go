package main

import (
	"fmt"
	"os"
	// "net"
	// "strconv"
	// "sync"
	// "time"
	// "regexp"
)

func main() {
	fmt.Println("test Network Vulnerability Scanner")

	root := "."
	scanDir := "./files"

	if !CheckGitignore(root) {
		fmt.Println("WARNING: no .gitignore found in", root)
		fmt.Println("Refusing to scan — without a .gitignore, secrets in ignored files could still leak.")
		os.Exit(1)
	}

	gi, err := NewGitignore(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read .gitignore: %v\n", err)
		os.Exit(1)
	}

	rules := []Rule{NewSecretsRule()}
	findings := scan(scanDir, rules, gi, 1)

	for _, f := range findings {
		fmt.Printf("%s:%d [%s] %s\n", f.Filename, f.Line, f.RuleID, f.Message)
	}
}
