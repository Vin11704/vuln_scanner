package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewUnpinnedActionsRule(t *testing.T) {
	rule := NewUnpinnedActionsRule()

	if rule == nil {
		t.Fatal("NewUnpinnedActionsRule() returned nil")
	}
	if rule.usesRegex == nil {
		t.Error("NewUnpinnedActionsRule() returned a rule with nil usesRegex")
	}
	if rule.shaRegex == nil {
		t.Error("NewUnpinnedActionsRule() returned a rule with nil shaRegex")
	}
}

func TestUnpinnedActionsRule_ID(t *testing.T) {
	rule := NewUnpinnedActionsRule()
	want := "github-unpinned-actions-scan"

	if got := rule.ID(); got != want {
		t.Errorf("ID() = %q, want %q", got, want)
	}
}

func TestUnpinnedActionsRule_AppliesTo(t *testing.T) {
	rule := NewUnpinnedActionsRule()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "yaml workflow", path: ".github/workflows/build.yml", want: true},
		{name: "yml workflow", path: ".github/workflows/build.yaml", want: true},
		{name: "nested workflow", path: "repo/.github/workflows/release/deploy.yml", want: true},
		{name: "windows workflow path", path: `repo\.github\workflows\build.yml`, want: true},
		{name: "outside workflows", path: ".github/actions/build.yml", want: false},
		{name: "wrong extension", path: ".github/workflows/build.json", want: false},
		{name: "uppercase extension", path: ".github/workflows/build.YML", want: true},
		{name: "empty path", path: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.AppliesTo(tt.path); got != tt.want {
				t.Errorf("AppliesTo(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestUnpinnedActionsRule_Check(t *testing.T) {
	rule := NewUnpinnedActionsRule()
	validSHA := strings.Repeat("a", 40)

	tests := []struct {
		name         string
		content      string
		wantFindings int
	}{
		{
			name:         "tag reference",
			content:      "uses: actions/checkout@v4",
			wantFindings: 1,
		},
		{
			name:         "branch reference",
			content:      "uses: actions/checkout@main",
			wantFindings: 1,
		},
		{
			name:         "expression reference",
			content:      "uses: actions/checkout@${{ inputs.ref }}",
			wantFindings: 1,
		},
		{
			name:         "short SHA reference",
			content:      "uses: actions/checkout@abc1234",
			wantFindings: 1,
		},
		{
			name:         "full lowercase SHA reference",
			content:      "uses: actions/checkout@" + validSHA,
			wantFindings: 0,
		},
		{
			name:         "indented comment",
			content:      "  # uses: actions/checkout@v4",
			wantFindings: 0,
		},
		{
			name:         "no uses reference",
			content:      "name: build\non: push",
			wantFindings: 0,
		},
		{
			name:         "empty content",
			content:      "",
			wantFindings: 0,
		},
		{
			name:         "multiple matching lines",
			content:      "steps:\n  - uses: actions/checkout@v4\n  - uses: actions/setup-go@main",
			wantFindings: 2,
		},
		{
			name:         "quoted reference is not matched",
			content:      `uses: "actions/checkout@` + validSHA + `"`,
			wantFindings: 0,
		},
		{
			name:         "only first reference on a line is checked",
			content:      "uses: actions/checkout@v4 uses: actions/setup-go@main",
			wantFindings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := rule.Check(".github/workflows/build.yml", []byte(tt.content))
			if len(findings) != tt.wantFindings {
				t.Errorf("Check() returned %d findings, want %d", len(findings), tt.wantFindings)
			}
		})
	}
}

func TestUnpinnedActionsRule_CheckFindingFields(t *testing.T) {
	rule := NewUnpinnedActionsRule()
	path := ".github/workflows/release.yml"
	content := "name: release\nsteps:\n  - uses: actions/checkout@v4"

	findings := rule.Check(path, []byte(content))
	if len(findings) != 1 {
		t.Fatalf("Check() returned %d findings, want 1", len(findings))
	}

	finding := findings[0]
	if finding.RuleID != rule.ID() {
		t.Errorf("RuleID = %q, want %q", finding.RuleID, rule.ID())
	}
	if finding.Filename != path {
		t.Errorf("Filename = %q, want %q", finding.Filename, path)
	}
	if finding.Line != 3 {
		t.Errorf("Line = %d, want 3", finding.Line)
	}
	if finding.Severity != "medium" {
		t.Errorf("Severity = %q, want %q", finding.Severity, "medium")
	}
	if !strings.Contains(finding.Message, "actions/checkout") || !strings.Contains(finding.Message, `"v4"`) {
		t.Errorf("Message = %q, want action and reference", finding.Message)
	}
}

func TestScanUnpinnedActionsRule(t *testing.T) {
	tempDir := t.TempDir()
	workflowDir := filepath.Join(tempDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	workflowPath := filepath.Join(workflowDir, "pipeline.yml")
	workflowContent := "steps:\n  - uses: actions/checkout@" + strings.Repeat("b", 40) + "\n  - uses: actions/setup-go@v5\n"
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		t.Fatalf("failed to create workflow file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "notes.txt"), []byte("uses: actions/checkout@v4\n"), 0644); err != nil {
		t.Fatalf("failed to create unrelated file: %v", err)
	}

	gitignore, err := NewGitignore(tempDir)
	if err != nil {
		t.Fatalf("NewGitignore failed: %v", err)
	}

	findings := scan(tempDir, []Rule{NewUnpinnedActionsRule()}, gitignore, 1)
	if len(findings) != 1 {
		t.Fatalf("scan returned %d findings, want 1", len(findings))
	}

	finding := findings[0]
	if finding.Filename != workflowPath {
		t.Errorf("Filename = %q, want %q", finding.Filename, workflowPath)
	}
	if finding.Line != 3 {
		t.Errorf("Line = %d, want 3", finding.Line)
	}
	if finding.RuleID != "github-unpinned-actions-scan" {
		t.Errorf("RuleID = %q, want %q", finding.RuleID, "github-unpinned-actions-scan")
	}
}
