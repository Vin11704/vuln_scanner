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

func NewUnpinnedActionsRule() *UnpinnedActionsRule{
    return &UnpinnedActionsRule{
        usesRegex: regexp.MustCompile(`uses:\s*([\w.-]+/[\w./-]+)@(\S+)`),
        shaRegex: regexp.MustCompile(`^[0-9a-f]{40}$`),
    }
}