package osexiterroranalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	identName = "os"
	exprName  = "Exit"
)

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
				if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := fun.X.(*ast.Ident); ok {
						if identName == ident.Name && exprName == fun.Sel.Name {
							pass.Reportf(x.Pos(), "os exit call error")
						}
					}
				}
			}

			return true
		})
	}
	return nil, nil
}
