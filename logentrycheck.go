package logentrycheck

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer"
)

func init() {
	register.Plugin("logentrycheck", New)
}

// Settings holds the plugin configuration.
type Settings struct{}

// Plugin implements register.LinterPlugin.
type Plugin struct{}

var _ register.LinterPlugin = &Plugin{}

// New creates a new plugin instance.
func New(_ any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

// BuildAnalyzers returns the list of analyzers.
func (*Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

// GetLoadMode returns the load mode for the plugin.
func (*Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
