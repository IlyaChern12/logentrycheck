package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/IlyaChern12/logentrycheck/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
