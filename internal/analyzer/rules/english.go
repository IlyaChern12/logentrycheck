package rules

import (
	"go/ast"
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// EnglishAnalyzer checks that log messages are in English only.
var EnglishAnalyzer = &analysis.Analyzer{
	Name:     "logentrycheck_english",
	Doc:      "checks that log messages are written in English only",
	Run:      runEnglish,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runEnglish(pass *analysis.Pass) (any, error) {
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

		checkEnglish(pass, msg, pos)
	})

	return nil, nil
}

// checkEnglish checks that a log message is in English only.
func checkEnglish(r reporter, msg string, pos token.Pos) {
	if !isEnglishOnly(msg) {
		r.Reportf(pos, "log message must be in English only: %q", msg)
	}
}

// isEnglishOnly checks that string contains only latin words.
func isEnglishOnly(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}

	return true
}
