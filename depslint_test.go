package depslint_test

import (
	"depslint"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	// modfile := testutil.ModFile(t, "./testdata/freee-transfer", nil)
	// testdata := testutil.WithModules(t, "./testdata/freee-transfer", modfile)
	// analysistest.Run(t, testdata, depslint.Analyzer, "./...")
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, depslint.Analyzer, "./...")

}
