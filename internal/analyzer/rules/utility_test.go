package rules_test

import (
	"go/ast"
	"testing"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func TestIsContextMethod(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "DebugContext",
			input: "DebugContext",
			want:  true,
		},
		{
			name:  "InfoContext",
			input: "InfoContext",
			want:  true,
		},
		{
			name:  "WarnContext",
			input: "WarnContext",
			want:  true,
		},
		{
			name:  "ErrorContext",
			input: "ErrorContext",
			want:  true,
		},
		{
			name:  "Debug",
			input: "Debug",
			want:  false,
		},
		{
			name:  "Info",
			input: "Info",
			want:  false,
		},
		{
			name:  "Warn",
			input: "Warn",
			want:  false,
		},
		{
			name:  "Error",
			input: "Error",
			want:  false,
		},
		{
			name:  "Fatal",
			input: "Fatal",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rules.IsContextMethod(tt.input); got != tt.want {
				t.Errorf("IsContextMethod(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsLoggerCall_Fallback(t *testing.T) {
	tests := []struct {
		name string
		call *ast.CallExpr
		want bool
	}{
		{
			name: "slog.Info",
			call: makeCall("slog", "Info", "user logged in"),
			want: true,
		},
		{
			name: "zap.Debug",
			call: makeCall("zap", "Debug", "request started"),
			want: true,
		},
		{
			name: "Sugar.Infof",
			call: makeCall("Sugar", "Infof", "user logged in"),
			want: true,
		},
		{
			name: "unknown package",
			call: makeCall("logrus", "Info", "user logged in"),
			want: false,
		},
		{
			name: "non selector expr",
			call: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "Info"},
				Args: []ast.Expr{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.IsLoggerCall(nil, tt.call)

			if got != tt.want {
				t.Errorf("IsLoggerCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
