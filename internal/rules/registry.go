package rules

import "scantrix/internal/types"

var registry []types.Rule

// RegisterRules lets rule files register themselves
func RegisterRules(rules []types.Rule) {
	registry = append(registry, rules...)
}

// LoadAll returns all registered rules
func LoadAll() []types.Rule {
	return registry
}
