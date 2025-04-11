// internal/astscanner/matchers/xss.go
package matchers

import (
	"scantrix/internal/astscanner"
	"scantrix/internal/types"
	"scantrix/pkg/php-ast/ast"
)

func init() {
	astscanner.RegisterMatcher("matchXSS", MatchXSS)
}

func MatchXSS(node ast.Node, ctx *astscanner.Context) *types.Finding {
	echoStmt, ok := node.(*ast.EchoStatement)
	if !ok {
		return nil
	}

	// Check if user input is being echoed directly
	if v, ok := echoStmt.Expr.(*ast.Variable); ok {
		if ctx.TaintedVars[v.Name] {
			return &types.Finding{
				RuleID:   "XSS001",
				Title:    "Reflected Cross-Site Scripting (XSS)",
				Severity: "critical",
				Advice:   "Escape output to prevent reflected XSS.",
				File:     ctx.File,
				Line:     v.Pos().Line,
			}
		}
	}

	return nil
}
