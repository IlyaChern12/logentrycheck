package rules_test

import (
	"testing"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		wantReports int
	}{
		{
			name:        "uppercase",
			msg:         "Starting server",
			wantReports: 1,
		},
		{
			name:        "lowercase",
			msg:         "starting server",
			wantReports: 0,
		},
		{
			name:        "empty message",
			msg:         "",
			wantReports: 0,
		},
		{
			name:        "digit first",
			msg:         "123 server started",
			wantReports: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockReporter{}
			rules.CheckLowercase(r, tt.msg, 0)

			if len(r.reports) != tt.wantReports {
				t.Errorf("checkLowercase() reports = %d, want %d",
					len(r.reports), tt.wantReports)
			}
		})
	}
}
