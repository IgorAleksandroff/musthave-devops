package checkers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_osExitChecker(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор
	// к пакетам из папки testdata и проверяет ожидания
	analysistest.Run(t, analysistest.TestData(), OSexitChecker, "./...")
}
