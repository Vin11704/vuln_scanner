package main


import (
	"testing"
)

// TestNewSecretsRule verifies the constructor returns a valid, usable rule.
func TestNewSecretsRule(t *testing.T) {
	rule := NewSecretsRule()

	if rule == nil {
		t.Fatal("NewSecretsRule() returned nil")
	}

	if len(rule.Pattern) == 0 {
		t.Fatal("NewSecretsRule() returned rule with no patterns")
	}

	for i, pat := range rule.Pattern {
		if pat.Regex == nil {
			t.Errorf("Pattern[%d] (%s) has nil Regex", i, pat.Name)
		}
		if pat.Name == "" {
			t.Errorf("Pattern[%d] has empty Name", i)
		}
		if pat.Severity == "" {
			t.Errorf("Pattern[%d] (%s) has empty Severity", i, pat.Name)
		}
	}
}

// TestSecretsRule_ID verifies the rule returns the correct identifier.
func TestSecretsRule_ID(t *testing.T) {
	rule := NewSecretsRule()
	got := rule.ID()
	want := "secret-scan"
	if got != want {
		t.Errorf("ID() = %q, want %q", got, want)
	}
}

// TestSecretsRule_AppliesTo verifies the rule applies to all paths.
func TestSecretsRule_AppliesTo(t *testing.T) {
	rule := NewSecretsRule()

	paths := []string{"app.go", "config.json", "README.md", "", "/some/deep/path.txt"}
	for _, p := range paths {
		if !rule.AppliesTo(p) {
			t.Errorf("AppliesTo(%q) = false, want true", p)
		}
	}
}

// TestSecretsRule_Check is a table-driven test covering all code paths of Check.
func TestSecretsRule_Check(t *testing.T) {
	rule := NewSecretsRule()

	tests := []struct {
		name         string
		path         string
		content      string
		wantFindings int
	}{
		// --- Non-JSON path ---
		{
			name:         "Non-JSON: no match",
			path:         "app.go",
			content:      "package main\nfunc main() {}\n",
			wantFindings: 0,
		},
		{
			name:         "Non-JSON: RSA private key found",
			path:         "config.txt",
			content:      "-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----",
			wantFindings: 1,
		},
		{
			name:         "Non-JSON: private key without RSA prefix",
			path:         "deploy.sh",
			content:      "#!/bin/bash\n-----BEGIN PRIVATE KEY-----\nMIIE...",
			wantFindings: 1,
		},
		{
			name:         "Non-JSON: key on line 3",
			path:         "file.txt",
			content:      "line one\nline two\n-----BEGIN RSA PRIVATE KEY-----\nline four",
			wantFindings: 1,
		},
		// --- JSON path ---
		{
			name:         "JSON: valid, no secrets",
			path:         "data.json",
			content:      `{"name":"safe","value":42}`,
			wantFindings: 0,
		},
		{
			name:         "JSON: top-level secret",
			path:         "config.json",
			content:      `{"key":"-----BEGIN PRIVATE KEY-----"}`,
			wantFindings: 1,
		},
		{
			name:         "JSON: nested object secret",
			path:         "nested.json",
			content:      `{"outer":{"inner":{"deep":"-----BEGIN RSA PRIVATE KEY-----"}}}`,
			wantFindings: 1,
		},
		{
			name:         "JSON: nested array secret",
			path:         "arr.json",
			content:      `{"items":["safe",["-----BEGIN PRIVATE KEY-----"]]}`,
			wantFindings: 1,
		},
		{
			name:         "JSON: invalid JSON",
			path:         "bad.json",
			content:      "not-valid-json{{{",
			wantFindings: 0,
		},
		{
			name:         "JSON: empty object",
			path:         "empty.json",
			content:      `{}`,
			wantFindings: 0,
		},
		// --- Edge cases ---
		{
			name:         "Empty content, non-JSON",
			path:         "empty.txt",
			content:      "",
			wantFindings: 0,
		},
		{
			name:         "Empty content, JSON extension",
			path:         "empty2.json",
			content:      "",
			wantFindings: 0,
		},
		{
			name:         "Case-insensitive JSON extension",
			path:         "Config.JSON",
			content:      `{"k":"-----BEGIN PRIVATE KEY-----"}`,
			wantFindings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := rule.Check(tt.path, []byte(tt.content))
			if len(findings) != tt.wantFindings {
				t.Errorf("Check(%q) returned %d findings, want %d", tt.path, len(findings), tt.wantFindings)
				for i, f := range findings {
					t.Logf("  finding[%d]: RuleID=%q File=%q Line=%d Sev=%q Msg=%q",
						i, f.RuleID, f.Filename, f.Line, f.Severity, f.Message)
				}
			}
		})
	}
}

// TestCheckFindingFields verifies that Finding struct fields are populated correctly.
func TestCheckFindingFields(t *testing.T) {
	rule := NewSecretsRule()

	t.Run("non-JSON finding fields", func(t *testing.T) {
		content := "line one\n-----BEGIN RSA PRIVATE KEY-----\nline three"
		findings := rule.Check("secret.txt", []byte(content))

		if len(findings) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(findings))
		}

		f := findings[0]
		if f.RuleID == "" {
			t.Error("RuleID is empty")
		}
		if f.Filename != "secret.txt" {
			t.Errorf("Filename = %q, want %q", f.Filename, "secret.txt")
		}
		if f.Line != 2 {
			t.Errorf("Line = %d, want 2", f.Line)
		}
		if f.Severity != "critical" {
			t.Errorf("Severity = %q, want %q", f.Severity, "critical")
		}
		if f.Message == "" {
			t.Error("Message is empty")
		}
	})

	t.Run("JSON finding fields", func(t *testing.T) {
		content := `{"secret":"-----BEGIN PRIVATE KEY-----"}`
		findings := rule.Check("creds.json", []byte(content))

		if len(findings) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(findings))
		}

		f := findings[0]
		if f.RuleID == "" {
			t.Error("RuleID is empty")
		}
		if f.Filename != "creds.json" {
			t.Errorf("Filename = %q, want %q", f.Filename, "creds.json")
		}
		if f.Severity != "critical" {
			t.Errorf("Severity = %q, want %q", f.Severity, "critical")
		}
		if f.Message == "" {
			t.Error("Message is empty")
		}
	})

	t.Run("non-JSON key on line 3", func(t *testing.T) {
		content := "aaa\nbbb\n-----BEGIN RSA PRIVATE KEY-----\nccc"
		findings := rule.Check("file.txt", []byte(content))

		if len(findings) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(findings))
		}
		if findings[0].Line != 3 {
			t.Errorf("Line = %d, want 3", findings[0].Line)
		}
	})
}

