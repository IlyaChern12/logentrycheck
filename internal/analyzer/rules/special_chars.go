package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// SpecialCharsAnalyzer checks that log messages don't contain special characters or emoji.
var SpecialCharsAnalyzer = &analysis.Analyzer{
	Name:     "logentrycheck_special_chars",
	Doc:      "checks that log messages don't contain special characters or emoji",
	Run:      runSpecialChars,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// allowedSpecialChars contains punctuation characters allowed in log messages.
var allowedSpecialChars = map[rune]bool{
	' ': true,
	'-': true,
	'_': true,
	'/': true,
	'.': true,
	',': true,
	'(': true,
	')': true,
	'[': true,
	']': true,
	'{': true,
	'}': true,
	'@': true,
	'#': true,
	'%': true,
	'+': true,
	'=': true,
	'<': true,
	'>': true,
	':': true,
}

func runSpecialChars(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		msg, pos, found := extractLogMessage(pass, call)
		if !found || msg == "" {
			return
		}

		checkSpecialChars(pass, msg, pos)
	})

	return nil, nil
}

// hasSpecialChars checks that string contains special characters or emoji.
func hasSpecialChars(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}

		if allowedSpecialChars[r] {
			continue
		}

		return true
	}

	return false
}

// checkSpecialChars checks that a log message does not contain special characters or emoji.
func checkSpecialChars(r reporter, msg string, pos token.Pos) {
	if !hasSpecialChars(msg) {
		return
	}

	fixed := removeSpecialChars(msg)

	if pass, ok := r.(*analysis.Pass); ok {
		pass.Report(analysis.Diagnostic{
			Pos:     pos,
			Message: fmt.Sprintf("log message must not contain special characters or emoji: %q", msg),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "remove special characters",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     pos,
							End:     pos + token.Pos(len(msg)+quoteLen),
							NewText: []byte(`"` + fixed + `"`),
						},
					},
				},
			},
		})
		return
	}

	r.Reportf(pos, "log message must not contain special characters or emoji: %q", msg)
}

// removeSpecialChars removes all disallowed characters from a string.
func removeSpecialChars(s string) string {
	var b strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || allowedSpecialChars[r] {
			b.WriteRune(r)
		}
	}

	return b.String()
}
