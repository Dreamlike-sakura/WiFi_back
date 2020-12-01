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
	app.POST("/verify", user.VerifyCodeHandler)
	app.POST("/check_user_info", user.UserInfoHandler)
	app.POST("/change_user_pwd", user.ChangePwdHandler)
	app.POST("/change_user_info", user.ChangeInfoHandler)
	app.POST("/check_user_movementlist", user.MovementListHandler)
	app.POST("/check_user_movement", user.MovementAmpPhaseHandler)
}
