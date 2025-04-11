// internal/astscanner/engine.go
package astscanner

import (
	"scantrix/pkg/php-ast/ast"
	"scantrix/internal/types"
)

type RuleFunc func(node ast.Node, ctx *Context) *types.Finding

type Context struct {
	File        string
	TaintedVars map[string]bool
}

var matcherRegistry = map[string]RuleFunc{}

func RegisterMatcher(name string, fn RuleFunc) {
	matcherRegistry[name] = fn
}

func RegisteredMatchers() map[string]RuleFunc {
	return matcherRegistry
}

func markTainted(node ast.Node, ctx *Context) {
	assign, ok := node.(*ast.Assignment)
	if !ok {
		return
	}

	if v, ok := assign.Value.(*ast.Variable); ok {
		if isSuperglobal(v.Name) {
			if name, ok := assign.Var.(*ast.Variable); ok {
				ctx.TaintedVars[name.Name] = true
			}
		}
	}
}

func isSuperglobal(name string) bool {
	switch name {
	case "$_GET", "$_POST", "$_REQUEST", "$_COOKIE", "$_SERVER":
		return true
	default:
		return false
	}
}

func Scan(nodes []ast.Node, rules []RuleFunc, filePath string) []types.Finding {
	var findings []types.Finding
	ctx := &Context{File: filePath, TaintedVars: map[string]bool{}}

	for _, node := range nodes {
		markTainted(node, ctx)
		for _, rule := range rules {
			if f := rule(node, ctx); f != nil {
				findings = append(findings, *f)
			}
		}
	}
	return findings
}
