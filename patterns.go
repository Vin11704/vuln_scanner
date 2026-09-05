package main
import(
	"regexp"
)

func NewSecretsRule() *SecretsRule{
	return &SecretsRule{
		Pattern: []secretPattern{
            {
                Name:     "private-key",
                Regex:    regexp.MustCompile(`-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`),
                Severity: "critical",
            },
        },
	}
}