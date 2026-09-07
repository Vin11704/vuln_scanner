package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	tempDir := t.TempDir()

	// Create a dummy .gitignore in tempDir
	giPath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(giPath, []byte("ignored.txt\n"), 0644); err != nil {
		t.Fatalf("failed to create dummy .gitignore: %v", err)
	}

	gi, err := NewGitignore(tempDir)
	if err != nil {
		t.Fatalf("NewGitignore failed: %v", err)
	}

	// Create a file containing a secret
	secretFile := filepath.Join(tempDir, "secret.txt")
	secretContent := []byte("api_key_data\n-----BEGIN RSA PRIVATE KEY-----\nsecret_key")
	if err := os.WriteFile(secretFile, secretContent, 0644); err != nil {
		t.Fatalf("failed to create secret file: %v", err)
	}

	// Create a clean file
	cleanFile := filepath.Join(tempDir, "clean.txt")
	cleanContent := []byte("hello world")
	if err := os.WriteFile(cleanFile, cleanContent, 0644); err != nil {
		t.Fatalf("failed to create clean file: %v", err)
	}

	rules := []Rule{NewSecretsRule()}

	// Test with 1 worker and 3 workers to test worker pool concurrency
	for _, workers := range []int{1, 3} {
		findings := scan(tempDir, rules, gi, workers)
		if len(findings) != 1 {
			t.Errorf("scan with %d workers returned %d findings, want 1", workers, len(findings))
		} else if findings[0].RuleID != "private-key" {
			t.Errorf("scan finding RuleID = %q, want 'private-key'", findings[0].RuleID)
		}
	}
}

