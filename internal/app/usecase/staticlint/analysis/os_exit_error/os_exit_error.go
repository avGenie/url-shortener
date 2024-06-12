// Package osexiterroranalyzer implements os.Exit constructions checker
package osexiterroranalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	identName = "os"
	exprName  = "Exit"
)

// OSExitAnalyzer os.Exit analyzer variable
var OSExitAnalyzer = &analysis.Analyzer{
	Name: "os_exit_error",
	Doc:  "check for os.Exit calls",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				fun, ok := x.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := fun.X.(*ast.Ident)
				if !ok {
					return true
				}

				if identName == ident.Name && exprName == fun.Sel.Name {
					pass.Reportf(x.Pos(), "os exit call error")
				}
			}

			return true
		})
	}
	return nil, nil
}
