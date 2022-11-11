package main

import (
	"github.com/ogugu9/depslint"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(depslint.Analyzer)
}
