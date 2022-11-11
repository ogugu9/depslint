package a

import (
	"fmt"
	// nolint:depslint
	"log" // aaa

	"go.uber.org/zap"
)

func f() {
	// The pattern can be written in regular expression.
	var gopher int
	print(gopher)
	fmt.Println("aa")
	log.Println("aa")
	logger, _ := zap.NewDevelopment()
	logger.Info("aaaa")
}
