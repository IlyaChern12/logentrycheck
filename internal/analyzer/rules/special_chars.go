package rules

import (
	"go/ast"
	"go/token"
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

		msg, pos, found := ExtractLogMessage(pass, call)
		if !found || msg == "" {
			return
		}

		CheckSpecialChars(pass, msg, pos)
	})

	return nil, nil
}

// HasSpecialChars checks that string contains special characters or emoji.
func HasSpecialChars(s string) bool {
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

// CheckSpecialChars checks that a log message does not contain special characters or emoji.
func CheckSpecialChars(r reporter, msg string, pos token.Pos) {
	if HasSpecialChars(msg) {
		r.Reportf(pos, "log message must not contain special characters or emoji: %q", msg)
	}
}
