package rules

import (
	"regexp"
	"scantrix/internal/types"
)

func CSRFRules() []types.Rule {
	return []types.Rule{
		{
			ID:       "CSRF001",
			Severity: "warning",
			FileTypes: []string{".php"},
			Pattern: regexp.MustCompile(`(?i)<form[^>]*method=['"]post['"][^>]*>`),
			Title:   "Possible CSRF (no CSRF token)",
			Advice:  "Ensure CSRF tokens are added to POST forms.",
		},
	}
}

func init() {
	RegisterRules(CSRFRules())
}