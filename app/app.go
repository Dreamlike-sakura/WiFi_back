package app

import (
	"github.com/gin-gonic/gin"
	"back/app/config"
	"back/app/model"
)

/**
 * 得到 Gin 实例
 */
func CreateGinApp() (app *gin.Engine) {
	// 得到参数
	config.Config.GetConfig()

	// 设置日志
	config.Init()

	app = gin.Default()

	// 链接数据库
	model.New()

	// 注册路由
	RegisterRouters(app)

	return
}


