// Package staticlint analyzes code using standard, staticcheck and self-designed analyzers
package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/avGenie/url-shortener/internal/app/usecase/staticlint/analysis"
	"github.com/avGenie/url-shortener/internal/app/usecase/staticlint/staticcheck"
)

func main() {
	analyzers := analysis.GetAnalyzers()
	analyzers = append(analyzers, staticcheck.GetAnalyzers()...)

	multichecker.Main(
		analyzers...,
	)
}
