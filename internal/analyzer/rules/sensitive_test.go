package rules_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func TestCheckSensitive(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		wantReports int
	}{
		{name: "clean message", msg: "user logged in", wantReports: 0},
		{name: "contains password", msg: "user password updated", wantReports: 0},
		{name: "contains token", msg: "invalid token received", wantReports: 0},
		{name: "contains api_key", msg: "api_key is missing", wantReports: 0},
		{name: "uppercase keyword", msg: "SECRET value exposed", wantReports: 0},
		{name: "token validated", msg: "token validated", wantReports: 0},
		{name: "user authenticated", msg: "user authenticated successfully", wantReports: 0},
		{name: "empty message", msg: "", wantReports: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockReporter{}
			msgArg := &ast.BasicLit{Kind: token.STRING, Value: `"` + tt.msg + `"`}

			rules.CheckSensitive(r, msgArg)

			if len(r.reports) != tt.wantReports {
				t.Errorf("CheckSensitive() reports = %d, want %d", len(r.reports), tt.wantReports)
			}
		})
	}
}

func TestCheckSensitive_NonSecure(t *testing.T) {
	tests := []struct {
		name        string
		call        *ast.CallExpr
		wantReports int
	}{
		{
			name: "password in concatenation",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "Info"},
				},
				Args: []ast.Expr{
					&ast.BinaryExpr{
						Op: token.ADD,
						X:  &ast.BasicLit{Kind: token.STRING, Value: `"user password: "`},
						Y:  &ast.Ident{Name: "password"},
					},
				},
			},
			wantReports: 1,
		},
		{
			name: "clean concatenation",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "Info"},
				},
				Args: []ast.Expr{
					&ast.BinaryExpr{
						Op: token.ADD,
						X:  &ast.BasicLit{Kind: token.STRING, Value: `"user id: "`},
						Y:  &ast.Ident{Name: "userID"},
					},
				},
			},
			wantReports: 0,
		},
		{
			name: "token in concatenation",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "Info"},
				},
				Args: []ast.Expr{
					&ast.BinaryExpr{
						Op: token.ADD,
						X:  &ast.BasicLit{Kind: token.STRING, Value: `"token: "`},
						Y:  &ast.Ident{Name: "tok"},
					},
				},
			},
			wantReports: 1,
		},
		{
			name: "nested concatenation with secret",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "Info"},
				},
				Args: []ast.Expr{
					&ast.BinaryExpr{
						Op: token.ADD,
						X: &ast.BinaryExpr{
							Op: token.ADD,
							X:  &ast.BasicLit{Kind: token.STRING, Value: `"user secret: "`},
							Y:  &ast.Ident{Name: "s"},
						},
						Y: &ast.BasicLit{Kind: token.STRING, Value: `" found"`},
					},
				},
			},
			wantReports: 1,
		},
		{
			name: "zap clean concatenation",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "zap"},
					Sel: &ast.Ident{Name: "Info"},
				},
				Args: []ast.Expr{
					&ast.BinaryExpr{
						Op: token.ADD,
						X:  &ast.BasicLit{Kind: token.STRING, Value: `"request id: "`},
						Y:  &ast.Ident{Name: "reqID"},
					},
				},
			},
			wantReports: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockReporter{}

			rules.CheckSensitive(r, tt.call.Args[0])

			if len(r.reports) != tt.wantReports {
				t.Errorf("CheckSensitive() reports = %d, want %d", len(r.reports), tt.wantReports)
			}
		})
	}
}

func TestContainsSensitiveKeyword(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantKeyword string
		wantFound   bool
	}{
		{
			name:        "clean message",
			input:       "user logged in",
			wantKeyword: "",
			wantFound:   false,
		},
		{
			name:        "password keyword",
			input:       "user password updated",
			wantKeyword: "password",
			wantFound:   true,
		},
		{
			name:        "token keyword",
			input:       "invalid token received",
			wantKeyword: "token",
			wantFound:   true,
		},
		{
			name:        "api key keyword",
			input:       "api_key is missing",
			wantKeyword: "api_key",
			wantFound:   true,
		},
		{
			name:        "uppercase keyword",
			input:       "SECRET value",
			wantKeyword: "secret",
			wantFound:   true,
		},
		{
			name:        "session keyword",
			input:       "session expired",
			wantKeyword: "session",
			wantFound:   true,
		},
		{
			name:        "credential keyword",
			input:       "credentials invalid",
			wantKeyword: "credential",
			wantFound:   true,
		},
		{
			name:        "empty string",
			input:       "",
			wantKeyword: "",
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeyword, gotFound := rules.ContainsSensitiveKeyword(tt.input)

			if gotFound != tt.wantFound || gotKeyword != tt.wantKeyword {
				t.Errorf("ContainsSensitiveKeyword(%q) = (%q, %v), want (%q, %v)",
					tt.input, gotKeyword, gotFound, tt.wantKeyword, tt.wantFound)
			}
		})
	}
}

func TestContainsSensitiveKeyword_CustomKeywords(t *testing.T) {
	tests := []struct {
		name        string
		keywords    string
		input       string
		wantKeyword string
		wantFound   bool
	}{
		{
			name:        "custom keyword found",
			keywords:    "mytoken,internalkey",
			input:       "mytoken is missing",
			wantKeyword: "mytoken",
			wantFound:   true,
		},
		{
			name:        "default keyword not found when custom set",
			keywords:    "mytoken",
			input:       "user password updated",
			wantKeyword: "",
			wantFound:   false,
		},
		{
			name:        "custom keyword with spaces",
			keywords:    " mytoken , internalkey ",
			input:       "internalkey exposed",
			wantKeyword: "internalkey",
			wantFound:   true,
		},
		{
			name:        "empty custom keywords uses defaults",
			keywords:    "",
			input:       "user password updated",
			wantKeyword: "password",
			wantFound:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*rules.CustomKeywordsFlag = tt.keywords
			defer func() { *rules.CustomKeywordsFlag = "" }()

			gotKeyword, gotFound := rules.ContainsSensitiveKeyword(tt.input)

			if gotFound != tt.wantFound || gotKeyword != tt.wantKeyword {
				t.Errorf("ContainsSensitiveKeyword(%q) = (%q, %v), want (%q, %v)",
					tt.input, gotKeyword, gotFound, tt.wantKeyword, tt.wantFound)
			}
		})
	}
}
