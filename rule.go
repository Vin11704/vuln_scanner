package main
import (
	"regexp"
	"strings"
	// "os"
	"encoding/json"
	"fmt"
)
type Finding struct{
	RuleID string
	Filename string
	Severity string
	Line int // which line of the file
	Message string // explanation of the finding
}

type Rule interface{
	ID() string
	AppliesTo(path string) bool
	Check(path string, content []byte) []Finding
}

type SecretsRule struct{
	Pattern []secretPattern
}

type secretPattern struct{
	Name string
	Regex *regexp.Regexp
	Severity string
}


func (r *SecretsRule) ID() string{
	return "secret-scan"
}

func (r *SecretsRule) AppliesTo(path string) bool{
	return true
}

func (r *SecretsRule) Check(path string, content []byte) []Finding{
	var findings []Finding
	if strings.HasSuffix(strings.ToLower(path), ".json") {

		var values map[string]interface{}
		err := json.Unmarshal(content, &values)
		
		if err != nil {
			return findings
		}
		
		// re := regexp.MustCompile(`-----BEGIN (RSA )?PRIVATE KEY-----`)

		// recursive check in case of nested json
		var walk func(v interface{}, pat secretPattern) bool
        walk = func(v interface{}, pat secretPattern) bool {
            switch x := v.(type) {
            case string:
                return pat.Regex.MatchString(x)
            case map[string]interface{}:
                for _, val := range x {
                    if walk(val, pat) {
                        return true
                    }
                }
            case []interface{}:
                for _, val := range x {
                    if walk(val, pat) {
                        return true
                    }
                }
            }
			return false	
			
		}

		for _, pat := range r.Pattern {
            if walk(values, pat) {
                findings = append(findings, Finding{
                    RuleID:   pat.Name,
                    Filename: path,
                    Severity: pat.Severity,
                    Message:  fmt.Sprintf("A %s was found in JSON file: %s", pat.Name, path),
                })
            }
        }
		
	return findings
	
	} else {
		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			for _, pat := range r.Pattern {
				if pat.Regex.MatchString(line) {
					findings = append(findings, Finding{
						RuleID:   pat.Name,
						Filename: path,
						Line:     lineNum + 1,
						Severity: pat.Severity,
						Message:  fmt.Sprintf("A %s was found in JSON file: %s", pat.Name, path),
					})
				}
			}
		}
		return findings
	}
}