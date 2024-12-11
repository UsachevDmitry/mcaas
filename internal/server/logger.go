package internal

import (
	"go.uber.org/zap"
)

var GlobalSugar zap.SugaredLogger

func Logger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	GlobalSugar = *logger.Sugar()
}
