package rules

import (
	"regexp"
	"scantrix/internal/types"
)

func XSSRules() []types.Rule {
	return []types.Rule{
		{
			ID:       "XSS001",
			Severity: "critical",
			FileTypes: []string{".js", ".html", ".twig"},
			Pattern: regexp.MustCompile(`(?i)document\.write\s*\(.*\$_(GET|POST|REQUEST)`),
			Title:   "Reflected XSS",
			Advice:  "Escape user input before rendering it into the DOM.",
		},
		{
			ID:       "XSS002",
			Severity: "warning",
			FileTypes: []string{".js"},
			Pattern: regexp.MustCompile(`innerHTML\s*=\s*`),
			Title:   "Possible DOM XSS",
			Advice:  "Avoid assigning raw HTML directly to innerHTML.",
		},
		{
			ID:       "XSS002",
			Severity: "info", // or "informational"
			FileTypes: []string{".js", ".php", ".html"},
			Pattern: regexp.MustCompile(`innerHTML\s*=\s*`),
			Title:   "Possible DOM XSS",
			Advice:  "Avoid assigning raw HTML directly to innerHTML.",
		},
	}
}

func init() {
	RegisterRules(XSSRules())
}