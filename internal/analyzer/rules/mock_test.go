package rules_test

import (
	"go/ast"
	"go/token"
)

type mockReporter struct {
	reports []mockReport
}

type mockReport struct {
	pos    token.Pos
	format string
	args   []any
}

func (m *mockReporter) Reportf(pos token.Pos, format string, args ...any) {
	m.reports = append(m.reports, mockReport{
		pos:    pos,
		format: format,
		args:   args,
	})
}

// makeCall creates a fake ast.CallExpr for pkg.Method("message").
func makeCall(pkg, method, msg string) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: pkg},
			Sel: &ast.Ident{Name: method},
		},
		Args: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"` + msg + `"`,
			},
		},
	}
}
