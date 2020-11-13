package controller

import (
	"back/app/config"
	"back/app/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	warp
}

func (s *User) LoginHandler(c *gin.Context) {
	user_name := c.Query("user_name")
	user_pwd := c.Query("user_pwd")
	if user_name == "" || user_pwd == "" {
		config.GetLogger().Warnw("账号密码不能为空",
			"user_name", user_name, "user_pwd", user_pwd,
		)
		c.JSON(http.StatusOK, s.FailWarp("账号密码不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetLoginData(user_name, user_pwd)
	if err != nil {
		config.GetLogger().Warnw("数据查询失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}
