// internal/astscanner/matchers/csrf.go
package matchers

import (
	"scantrix/internal/astscanner"
	"scantrix/internal/types"
	"scantrix/pkg/php-ast/ast"
	"strings"
)

func init() {
	astscanner.RegisterMatcher("matchCSRF", MatchCSRF)
}

func MatchCSRF(node ast.Node, ctx *astscanner.Context) *types.Finding {
	formOutput, ok := node.(*ast.EchoStatement)
	if !ok {
		return nil
	}

	// Naively detect HTML form outputs without CSRF protection
	if lit, ok := formOutput.Expr.(*ast.Literal); ok {
		lower := strings.ToLower(lit.Value)
		if strings.Contains(lower, "<form") &&
			strings.Contains(lower, "post") &&
			!strings.Contains(lower, "csrf") {
			return &types.Finding{
				RuleID:   "CSRF001",
				Title:    "Form missing CSRF protection",
				Severity: "warning",
				Advice:   "Include a CSRF token in all HTML forms.",
				File:     ctx.File,
				Line:     lit.Pos().Line,
			}
		}
	}

	return nil
}