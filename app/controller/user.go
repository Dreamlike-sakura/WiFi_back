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

func (s *User) RegisterHandler(c *gin.Context) {
	i := new(model.Info)
	i.User = c.Query("user_name")
	i.Password = c.Query("user_pwd")
	i.Email = c.Query("user_email")
	i.Tel = c.Query("user_tel")
	i.Sex = "M"
	i.Type = "0"
	i.Head_portrait = "1"

	if i.User == "" || i.Password == "" {
		config.GetLogger().Warnw("账户信息不能为空",
			"user_name:", i.User, "user_pwd:", i.Password,
		)
		c.JSON(http.StatusOK, s.FailWarp("账号信息不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetRegisterData(i)
	if err != nil {
		config.GetLogger().Warnw("注册数据查询失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}

func (s *User) SecureCodeHandler(c *gin.Context) {
	tel := c.Query("tel")

	if tel == "" {
		config.GetLogger().Warnw("手机号码不能为空",
			"tel:", tel,
		)
		c.JSON(http.StatusOK, s.FailWarp("手机号码不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetSecureCodeData(tel)
	if err != nil {
		config.GetLogger().Warnw("发送手机验证码失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}