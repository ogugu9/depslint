package bb

import (
	"log"

	"github.com/ogugu9/a/aa"
	//lint:ignore depslint
	"github.com/ogugu9/a/cc"
	"go.uber.org/zap"
)

func BBB() {
	aa.AAA()
	cc.CCC()
	log.Println("AAA")
	logger, _ := zap.NewDevelopment()
	logger.Info("AAA")
}
