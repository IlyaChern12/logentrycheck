package rules_test

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func testdata(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	return filepath.Join(wd, "testdata")
}

func TestLowercaseAnalyzer(t *testing.T) {
	analysistest.Run(t, testdata(t), rules.LowercaseAnalyzer, "lowercase")
}

func TestLowercaseAnalyzerFix(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, testdata(t), rules.LowercaseAnalyzer, "lowercase")
}

func TestEnglishAnalyzer(t *testing.T) {
	analysistest.Run(t, testdata(t), rules.EnglishAnalyzer, "english")
}

func TestSpecialCharsAnalyzer(t *testing.T) {
	analysistest.Run(t, testdata(t), rules.SpecialCharsAnalyzer, "special_chars")
}

func TestSpecialCharsAnalyzerFix(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, testdata(t), rules.SpecialCharsAnalyzer, "special_chars")
}

func TestSensitiveAnalyzer(t *testing.T) {
	analysistest.Run(t, testdata(t), rules.SensitiveAnalyzer, "sensitive")
}
