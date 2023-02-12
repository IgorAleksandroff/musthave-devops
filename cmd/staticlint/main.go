package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"

	"github.com/IgorAleksandroff/musthave-devops/internal/checkers"
)

func main() {
	mychecks := []*analysis.Analyzer{
		checkers.OSexitChecker,
		copylock.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		unusedresult.Analyzer,
	}

	for _, v := range simple.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	multichecker.Main(
		mychecks...,
	)
}
