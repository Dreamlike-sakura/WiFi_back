package app

import (
	"back/app/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(app *gin.Engine) {
	user := controller.User{}
	app.POST("/login", user.LoginHandler)
	app.POST("/register", user.RegisterHandler)
	app.POST("/send_code", user.SecureCodeHandler)
	app.POST("/find_password", user.VerifyCodeHandler)
	app.POST("/check_user_info", user.UserInfoHandler)
}
