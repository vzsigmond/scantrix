package rules

import (
	"regexp"
	"scantrix/internal/types"
)

func OpenRedirectRules() []types.Rule {
	return []types.Rule{
		{
			ID:       "REDIRECT001",
			Severity: "critical",
			FileTypes: []string{".php", ".js"},
			Pattern: regexp.MustCompile(`(?i)header\s*\(\s*['"]Location:.*\$_(GET|POST|REQUEST)\[`),
			Title:   "Possible Open Redirect (PHP)",
			Advice:  "Ensure redirects use a whitelist of trusted URLs.",
		},
		{
			ID:       "REDIRECT002",
			Severity: "warning",
			FileTypes: []string{".js"},
			Pattern: regexp.MustCompile(`(?i)location\.href\s*=\s*\w+\s*[\+\=]`),
			Title:   "JavaScript Open Redirect",
			Advice:  "Validate URL origin before redirecting.",
		},
	}
}

func init() {
	RegisterRules(OpenRedirectRules())
}