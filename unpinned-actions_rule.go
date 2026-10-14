package main
import(
	"regexp"
	"strings"
	"path/filepath"

	// "os"
	"fmt"
)

type UnpinnedActionsRule struct {
	usesRegex	*regexp.Regexp
	shaRegex	*regexp.Regexp
}

func (r *UnpinnedActionsRule) ID() string {
	return "github-unpinned-actions-scan"
}

func (r *UnpinnedActionsRule) AppliesTo(path string) bool {
	p := filepath.ToSlash(strings.ToLower(path))
	return strings.Contains(p, ".github/workflows/") &&
	(strings.HasSuffix(p, ".yml") || strings.HasSuffix(p, ".yaml"))
}

func (r *UnpinnedActionsRule) Check(path string, content []byte) []Finding{
	var findings []Finding
	for i, line := range strings.Split(string(content), "\n"){
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		matches := r.usesRegex.FindStringSubmatch(line)
		if matches==nil{continue}

		action, ref := matches[1], matches[2]
		if !r.shaRegex.MatchString(ref){
			findings = append(findings, Finding{
					RuleID:   r.ID(),
					Filename: path,
					Severity: "medium",
					Line: i+1,
					Message:  fmt.Sprintf("%s is pinned to %q, not a commit SHA", action, ref),
				})
		}
	}
	return findings
}