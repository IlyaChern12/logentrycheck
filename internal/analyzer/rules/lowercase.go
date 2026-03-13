package rules

import (
	"fmt"
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

const quoteLen = 2

// checkLowercase checks a single call expression against the lowercase rule.
func checkLowercase(r reporter, msg string, pos token.Pos) {
	if msg == "" || !unicode.IsUpper(rune(msg[0])) {
		return
	}

	fixed := string(unicode.ToLower(rune(msg[0]))) + msg[1:]

	if pass, ok := r.(*analysis.Pass); ok {
		pass.Report(analysis.Diagnostic{
			Pos:     pos,
			Message: fmt.Sprintf("log message should start with a lowercase letter: %q", msg),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "convert first letter to lowercase",
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

	r.Reportf(pos, "log message should start with a lowercase letter: %q", msg)
}
