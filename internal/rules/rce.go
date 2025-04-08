package rules

import (
	"regexp"
	"scantrix/internal/types"
)

func RCERules() []types.Rule {
	return []types.Rule{
		{
			ID:       "RCE001",
			Severity: "critical",
			FileTypes: []string{".php"},
			Pattern: regexp.MustCompile(`(?i)(eval|system|exec|shell_exec|passthru)\s*\(\s*\$_(GET|POST|REQUEST|SERVER)\[`),
			Title:   "Remote Code Execution via user input",
			Advice:  "Avoid using input directly in execution functions.",
		},
	}
}

func init() {
	RegisterRules(RCERules())
}