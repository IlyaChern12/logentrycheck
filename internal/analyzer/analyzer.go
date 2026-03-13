package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

// Analyzer is the entry point for the logentrycheck linter.
var Analyzer = &analysis.Analyzer{
	Name: "logentrycheck",
	Doc:  "checks log entries for common mistakes in log messages",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var (
	disableLowercase    bool
	disableEnglish      bool
	disableSpecialChars bool
	disableSensitive    bool
)

func init() {
	Analyzer.Flags.BoolVar(&disableLowercase, "disable-lowercase", false, "disable lowercase rule")
	Analyzer.Flags.BoolVar(&disableEnglish, "disable-english", false, "disable english rule")
	Analyzer.Flags.BoolVar(&disableSpecialChars, "disable-special-chars", false, "disable special chars rule")
	Analyzer.Flags.BoolVar(&disableSensitive, "disable-sensitive", false, "disable sensitive rule")
}

func run(pass *analysis.Pass) (any, error) {
	for _, rule := range rules.Active(disableLowercase, disableEnglish, disableSpecialChars, disableSensitive) {
		if _, err := rule.Run(pass); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
