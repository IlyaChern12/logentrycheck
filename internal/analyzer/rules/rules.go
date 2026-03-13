package rules

import "golang.org/x/tools/go/analysis"

// All contains registered rules.
var All = []*analysis.Analyzer{
	LowercaseAnalyzer,
	EnglishAnalyzer,
	SpecialCharsAnalyzer,
	SensitiveAnalyzer,
}

// Active returns the list of enabled analyzers based on disable flags.
func Active(disableLowercase, disableEnglish, disableSpecialChars, disableSensitive bool) []*analysis.Analyzer {
	var active []*analysis.Analyzer

	if !disableLowercase {
		active = append(active, LowercaseAnalyzer)
	}

	if !disableEnglish {
		active = append(active, EnglishAnalyzer)
	}

	if !disableSpecialChars {
		active = append(active, SpecialCharsAnalyzer)
	}

	if !disableSensitive {
		active = append(active, SensitiveAnalyzer)
	}

	return active
}
