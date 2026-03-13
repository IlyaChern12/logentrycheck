package logentrycheck

import (
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer"
	"github.com/IlyaChern12/logentrycheck/internal/analyzer/rules"
)

func init() {
	register.Plugin("logentrycheck", New)
}

// Settings holds the plugin configuration from .golangci.yml.
type Settings struct {
	DisableLowercase    bool     `json:"disableLowercase"`
	DisableEnglish      bool     `json:"disableEnglish"`
	DisableSpecialChars bool     `json:"disableSpecialChars"`
	DisableSensitive    bool     `json:"disableSensitive"`
	Keywords            []string `json:"keywords"`
}

// Plugin implements register.LinterPlugin.
type Plugin struct {
	settings Settings
}

var _ register.LinterPlugin = &Plugin{}

// New creates a new plugin instance with decoded settings.
func New(input any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](input)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: settings}, nil
}

// BuildAnalyzers returns the list of analyzers with applied settings.
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if err := analyzer.Analyzer.Flags.Set("disable-lowercase", boolToString(p.settings.DisableLowercase)); err != nil {
		return nil, err
	}

	if err := analyzer.Analyzer.Flags.Set("disable-english", boolToString(p.settings.DisableEnglish)); err != nil {
		return nil, err
	}

	if err := analyzer.Analyzer.Flags.Set("disable-special-chars", boolToString(p.settings.DisableSpecialChars)); err != nil {
		return nil, err
	}

	if err := analyzer.Analyzer.Flags.Set("disable-sensitive", boolToString(p.settings.DisableSensitive)); err != nil {
		return nil, err
	}

	if len(p.settings.Keywords) > 0 {
		if err := rules.SensitiveAnalyzer.Flags.Set("keywords", strings.Join(p.settings.Keywords, ",")); err != nil {
			return nil, err
		}
	}

	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

// GetLoadMode returns the load mode for the plugin.
func (*Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func boolToString(b bool) string {
	if b {
		return "true"
	}

	return "false"
}
