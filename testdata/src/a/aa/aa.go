package aa

import (
	"log"

	"go.uber.org/zap"
)

func AAA() {
	log.Println("AAA")
	logger, _ := zap.NewDevelopment()
	logger.Info("AAA")
}
