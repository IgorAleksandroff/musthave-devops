package checkers

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var OSexitChecker = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit in function main package main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if funcDecl, ok := node.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == "main" {
					ast.Inspect(file, func(node ast.Node) bool {
						if callExpr, ok := node.(*ast.CallExpr); ok {
							if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								name := selectorExpr.Sel.Name
								if strings.Contains("os.Exit", name) {
									pass.Reportf(selectorExpr.Pos(), "function main should not have os exit")
								}

							}
						}
						return true
					})
				}
			}
			return true
		})
	}
	return nil, nil
}
