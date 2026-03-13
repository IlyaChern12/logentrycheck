package rules

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// SensitiveAnalyzer checks that log messages don't contain potentially sensitive data.
var SensitiveAnalyzer = &analysis.Analyzer{
	Name:     "logentrycheck_sensitive",
	Doc:      "checks that log messages do not contain potentially sensitive data",
	Run:      runSensitive,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// customKeywordsFlag holds comma-separated custom sensitive keywords from flags.
var customKeywordsFlag string

func init() {
	SensitiveAnalyzer.Flags.StringVar(
		&customKeywordsFlag,
		"keywords",
		"",
		"comma-separated list of additional sensitive keywords",
	)
}

// defaultSensitiveKeywords contains default keywords indicating sensitive data.
var defaultSensitiveKeywords = []string{
	"password", "passwd", "secret", "token",
	"api_key", "apikey", "auth", "credential",
	"private_key", "access_key", "session",
}

func runSensitive(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		msgArg := extractMessageArg(pass, call)
		if msgArg == nil {
			return
		}

		checkSensitive(pass, msgArg)
	})

	return nil, nil
}

// сontainsSensitiveKeyword checks that the message contains a sensitive keyword.
func containsSensitiveKeyword(msg string) (string, bool) {
	lower := strings.ToLower(msg)

	keywords := defaultSensitiveKeywords
	if customKeywordsFlag != "" {
		keywords = strings.Split(customKeywordsFlag, ",")
		for i, k := range keywords {
			keywords[i] = strings.TrimSpace(k)
		}
	}

	for _, keyword := range keywords {
		if keyword != "" && strings.Contains(lower, keyword) {
			return keyword, true
		}
	}

	return "", false
}

// checkSensitive checks that a log message doesn't contain sensitive data.
func checkSensitive(r reporter, msgArg ast.Expr) {
	if _, ok := msgArg.(*ast.BinaryExpr); !ok {
		return
	}

	checkExprForSensitive(r, msgArg)
}

// checkExprForSensitive recursively checks all string literals in an expression.
func checkExprForSensitive(r reporter, expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return
		}

		lit := e.Value[1 : len(e.Value)-1]
		if keyword, found := containsSensitiveKeyword(lit); found {
			r.Reportf(e.Pos(), "log message may contain sensitive data (keyword: %q): %q", keyword, lit)
		}

	case *ast.BinaryExpr:
		checkExprForSensitive(r, e.X)
		checkExprForSensitive(r, e.Y)
	}
}
