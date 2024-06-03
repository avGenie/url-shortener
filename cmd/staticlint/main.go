// Package staticlint analyzes code using standard and self-designed analyzers
package main

import "github.com/avGenie/url-shortener/internal/app/usecase/staticlint/analysis"

func main() {
	analysis.Analyse()
}
