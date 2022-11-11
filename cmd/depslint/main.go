package main

import (
	"depslint"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(depslint.Analyzer)
}
