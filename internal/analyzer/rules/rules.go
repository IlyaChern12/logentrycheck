package rules

import "golang.org/x/tools/go/analysis"

// All contains registered rules.
var All = []*analysis.Analyzer{
	LowercaseAnalyzer,
	EnglishAnalyzer,
	SpecialCharsAnalyzer,
	SensitiveAnalyzer,
}
