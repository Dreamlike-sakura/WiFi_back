package config

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func Init() {
	var temp *zap.Logger

	temp, _ = zap.NewDevelopment()
	logger = temp.Sugar()

	logger.Info("[Init] 初始化日志成功")
}

func GetLogger() *zap.SugaredLogger {
	return logger
}
