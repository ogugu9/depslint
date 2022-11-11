package cc

import (
	"log"

	"go.uber.org/zap"
)

func CCC() {
	log.Println("CCc")
	logger, _ := zap.NewDevelopment()
	logger.Info("CCC")
}
