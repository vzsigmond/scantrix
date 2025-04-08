package rules

import (
	"regexp"
	"scantrix/internal/types"
)

func InsecureCryptoRules() []types.Rule {
	return []types.Rule{
		{
			ID:       "CRYPTO001",
			Severity: "warning",
			FileTypes: []string{".php", ".js"},
			Pattern: regexp.MustCompile(`(?i)md5\s*\(`),
			Title:   "Insecure Hash Function (MD5)",
			Advice:  "Avoid using MD5. Use stronger hashes like SHA-256 or bcrypt.",
		},
	}
}

func init() {
	RegisterRules(InsecureCryptoRules())
}
