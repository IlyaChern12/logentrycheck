package rules_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func TestExtractMessageArg(t *testing.T) {
	tests := []struct {
		name    string
		call    *ast.CallExpr
		wantNil bool
	}{
		{
			name:    "slog Info returns arg",
			call:    makeCall("slog", "Info", "user logged in"),
			wantNil: false,
		},
		{
			name:    "zap Infof returns arg",
			call:    makeCall("zap", "Infof", "request completed"),
			wantNil: false,
		},
		{
			name:    "unknown package returns nil",
			call:    makeCall("logrus", "Info", "user logged in"),
			wantNil: true,
		},
		{
			name:    "unknown method returns nil",
			call:    makeCall("slog", "Trace", "user logged in"),
			wantNil: true,
		},
		{
			name: "no args returns nil",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "Info"},
				},
				Args: []ast.Expr{},
			},
			wantNil: true,
		},
		{
			name: "non selector fun returns nil",
			call: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "Info"},
				Args: []ast.Expr{},
			},
			wantNil: true,
		},
		{
			name: "basic lit fun returns nil",
			call: &ast.CallExpr{
				Fun: &ast.BasicLit{Kind: token.STRING, Value: `"slog"`},
			},
			wantNil: true,
		},
		{
			name: "context method returns arg at index 1",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "InfoContext"},
				},
				Args: []ast.Expr{
					&ast.Ident{Name: "ctx"},
					&ast.BasicLit{Kind: token.STRING, Value: `"user logged in"`},
				},
			},
			wantNil: false,
		},
		{
			name: "context method with only ctx arg returns nil",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "InfoContext"},
				},
				Args: []ast.Expr{
					&ast.Ident{Name: "ctx"},
				},
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ExtractMessageArg(nil, tt.call)

			if tt.wantNil && got != nil {
				t.Errorf("ExtractMessageArg() = %v, want nil", got)
			}

			if !tt.wantNil && got == nil {
				t.Errorf("ExtractMessageArg() = nil, want non-nil")
			}
		})
	}
}
