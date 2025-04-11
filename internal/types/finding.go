// internal/types/finding.go
package types

type Finding struct {
	RuleID   string `json:"rule_id"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
	Advice   string `json:"advice"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}
