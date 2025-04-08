package rules

import (
	"regexp"
	"scantrix/internal/types"
)

func SQLInjectionRules() []types.Rule {
	return []types.Rule{
		{
			ID:       "SQLI001",
			Severity: "critical",
			FileTypes: []string{".php", ".js"},
			Pattern: regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*(\$_(GET|POST|REQUEST))`),
			Title:   "Possible SQL Injection",
			Advice:  "Use prepared statements or parameterized queries.",
		},
	}
}

func init() {
	RegisterRules(SQLInjectionRules())
}