// Package staticcheck works with staticcheck analysers from honnef.co/go/tools/staticcheck
package staticcheck

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

var (
	simpleAnalyzers = map[string]struct{}{
		"S1011": struct{}{},
		"S1017": struct{}{},
	}
	stylecheckAnalyzers = map[string]struct{}{"ST1023": struct{}{}}
	quickfixAnalyzers   = map[string]struct{}{"QF1007": struct{}{}}
)

// GetAnalyzers Returns all SA staticcheck analyzers
//
// List of analyzers: https://staticcheck.io/docs/checks
func GetAnalyzers() []*analysis.Analyzer {
	count := len(staticcheck.Analyzers) + len(simpleAnalyzers) +
		len(stylecheckAnalyzers) + len(quickfixAnalyzers)
	analyzers := make([]*analysis.Analyzer, 0, count)
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	for _, v := range simple.Analyzers {
		if _, ok := simpleAnalyzers[v.Analyzer.Name]; ok {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if _, ok := stylecheckAnalyzers[v.Analyzer.Name]; ok {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	for _, v := range quickfix.Analyzers {
		if _, ok := quickfixAnalyzers[v.Analyzer.Name]; ok {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	return analyzers
}
