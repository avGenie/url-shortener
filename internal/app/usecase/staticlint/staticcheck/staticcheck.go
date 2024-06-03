// Package staticcheck works with staticcheck analysers from honnef.co/go/tools/staticcheck
package staticcheck

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

// GetAnalyzers Returns all SA staticcheck analyzers
//
// List of analyzers: https://staticcheck.io/docs/checks
func GetAnalyzers() []*analysis.Analyzer {
	analyzers := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	return analyzers
}
