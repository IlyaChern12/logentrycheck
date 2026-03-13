package rules

import (
	"go/ast"
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// LowercaseAnalyzer checks that log messages start with a lowercase letter.
var LowercaseAnalyzer = &analysis.Analyzer{
	Name: "logentrycheck_lowercase",
	Doc:  "checks that log messages start with a lowercase letter",
	Run:  runLowercase,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runLowercase(pass *analysis.Pass) (any, error) {
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

		checkLowercase(pass, msg, pos)
	})

	return nil, nil
}

// checkLowercase checks a single call expression against the lowercase rule.
func checkLowercase(r reporter, msg string, pos token.Pos) {
	if msg == "" {
		return
	}

	if unicode.IsUpper(rune(msg[0])) {
		r.Reportf(pos, "log message should start with a lowercase letter: %q", msg)
	}
}
