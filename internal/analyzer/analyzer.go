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

func run(pass *analysis.Pass) (any, error) {
	for _, rule := range rules.All {
		_, err := rule.Run(pass)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
