package config

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func Init() {
	var temp *zap.Logger
	if Config.App.Mode == "release" {
		temp, _ = zap.NewProduction()
	} else {
		temp, _ = zap.NewDevelopment()
	}

	logger = temp.Sugar()

	logger.Info("[Init] 初始化日志成功")
}

func GetLogger() *zap.SugaredLogger {
	return logger
}
