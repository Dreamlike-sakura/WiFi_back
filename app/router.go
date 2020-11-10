package app

import (
	"back/app/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(app *gin.Engine) {
	user := controller.User{}
	app.GET("/login", user.LoginHandler)
}
