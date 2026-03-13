package rules_test

import (
	"testing"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		name        string
		pkg         string
		method      string
		msg         string
		wantReports int
	}{
		{
			name:        "ascii message no report",
			pkg:         "slog",
			method:      "Info",
			msg:         "user logged in",
			wantReports: 0,
		},
		{
			name:        "cyrillic message report",
			pkg:         "slog",
			method:      "Info",
			msg:         "пользователь вошёл",
			wantReports: 1,
		},
		{
			name:        "emoji report",
			pkg:         "slog",
			method:      "Info",
			msg:         "user logged in 🎉",
			wantReports: 1,
		},
		{
			name:        "mixed ascii and non ascii report",
			pkg:         "zap",
			method:      "Info",
			msg:         "user üser alles",
			wantReports: 1,
		},
		{
			name:        "empty message no report",
			pkg:         "slog",
			method:      "Info",
			msg:         "",
			wantReports: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockReporter{}

			rules.CheckEnglish(r, tt.msg, 0)

			if len(r.reports) != tt.wantReports {
				t.Errorf("CheckEnglish() reports = %d, want %d", len(r.reports), tt.wantReports)
			}
		})
	}
}

func TestIsEnglishOnly(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "ascii letters",
			input: "hello world",
			want:  true,
		},
		{
			name:  "ascii with digits and symbols",
			input: "abc123!@#",
			want:  true,
		},
		{
			name:  "empty string",
			input: "",
			want:  true,
		},
		{
			name:  "cyrillic",
			input: "пользователь",
			want:  false,
		},
		{
			name:  "emoji",
			input: "user 🎉",
			want:  false,
		},
		{
			name:  "latin extended",
			input: "üser",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rules.IsEnglishOnly(tt.input); got != tt.want {
				t.Errorf("IsEnglishOnly(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
