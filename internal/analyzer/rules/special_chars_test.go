package rules_test

import (
	"testing"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func TestCheckSpecialChars(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		wantReports int
	}{
		{
			name:        "plain text",
			msg:         "user logged in",
			wantReports: 0,
		},
		{
			name:        "allowed punctuation",
			msg:         "user/admin logged in (role=admin)",
			wantReports: 0,
		},
		{
			name:        "allowed special chars",
			msg:         "request [GET /api/v1] status=200",
			wantReports: 0,
		},
		{
			name:        "exclamation mark",
			msg:         "user logged in!",
			wantReports: 1,
		},
		{
			name:        "emoji",
			msg:         "user logged in 🎉",
			wantReports: 1,
		},
		{
			name:        "question mark",
			msg:         "what happened?",
			wantReports: 1,
		},
		{
			name:        "empty message",
			msg:         "",
			wantReports: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockReporter{}

			rules.CheckSpecialChars(r, tt.msg, 0)

			if len(r.reports) != tt.wantReports {
				t.Errorf("got %d reports, want %d", len(r.reports), tt.wantReports)
			}
		})
	}
}

func TestHasSpecialChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "plain text",
			input: "user logged in",
			want:  false,
		},
		{
			name:  "path and slashes",
			input: "request GET /api/v1",
			want:  false,
		},
		{
			name:  "parens and equals",
			input: "status=200 (ok)",
			want:  false,
		},
		{
			name:  "square brackets",
			input: "role [admin]",
			want:  false,
		},
		{
			name:  "plus sign",
			input: "key=value+extra",
			want:  false,
		},
		{
			name:  "exclamation mark",
			input: "user logged in!",
			want:  true,
		},
		{
			name:  "question mark",
			input: "what happened?",
			want:  true,
		},
		{
			name:  "emoji",
			input: "user 🎉",
			want:  true,
		},
		{
			name:  "cyrillic with special char",
			input: "привет!",
			want:  true,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := rules.HasSpecialChars(tt.input); got != tt.want {
				t.Errorf("HasSpecialChars(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
