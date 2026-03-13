package rules_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func TestActive(t *testing.T) {
	tests := []struct {
		name                string
		disableLowercase    bool
		disableEnglish      bool
		disableSpecialChars bool
		disableSensitive    bool
		expected            []*analysis.Analyzer
	}{
		{
			name: "all enabled",
			expected: []*analysis.Analyzer{
				rules.LowercaseAnalyzer,
				rules.EnglishAnalyzer,
				rules.SpecialCharsAnalyzer,
				rules.SensitiveAnalyzer,
			},
		},
		{
			name:             "disable lowercase",
			disableLowercase: true,
			expected: []*analysis.Analyzer{
				rules.EnglishAnalyzer,
				rules.SpecialCharsAnalyzer,
				rules.SensitiveAnalyzer,
			},
		},
		{
			name:                "disable english and special",
			disableEnglish:      true,
			disableSpecialChars: true,
			expected: []*analysis.Analyzer{
				rules.LowercaseAnalyzer,
				rules.SensitiveAnalyzer,
			},
		},
		{
			name:                "disable all",
			disableLowercase:    true,
			disableEnglish:      true,
			disableSpecialChars: true,
			disableSensitive:    true,
			expected:            []*analysis.Analyzer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rules.Active(
				tt.disableLowercase,
				tt.disableEnglish,
				tt.disableSpecialChars,
				tt.disableSensitive,
			)

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d analyzers, got %d", len(tt.expected), len(result))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Fatalf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
