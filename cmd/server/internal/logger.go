package internal

import (
	"go.uber.org/zap"
)

var GlobalSugar zap.SugaredLogger

func Logger(){
		// создаём предустановленный регистратор zap
		logger, err := zap.NewDevelopment()
		if err != nil {
			// вызываем панику, если ошибка
			panic(err)
		}
		defer logger.Sync()
	
		// делаем регистратор SugaredLogger
		GlobalSugar = *logger.Sugar()
	
}