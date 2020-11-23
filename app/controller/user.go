package controller

import (
	"back/app/config"
	"back/app/model"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type User struct {
	warp
}

func (s *User) LoginHandler(c *gin.Context) {
	//user_name := c.Query("user_name")
	//user_pwd := c.Query("user_pwd")
	cont, _ := ioutil.ReadAll(c.Request.Body)
	if cont == nil {
		config.GetLogger().Warnw("账号密码不能为空",
			"cont", cont,
		)
		c.JSON(http.StatusOK, s.FailWarp("账号密码不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetLoginData(string(cont))
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

	cont, _ := ioutil.ReadAll(c.Request.Body)

	if cont ==  nil {
		config.GetLogger().Warnw("账户信息不能为空",
			"cont:", cont,
		)
		c.JSON(http.StatusOK, s.FailWarp("账号信息不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetRegisterData(string(cont))
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
	tel, _ := ioutil.ReadAll(c.Request.Body)

	if tel == nil {
		config.GetLogger().Warnw("手机号码不能为空",
			"tel:", tel,
		)
		c.JSON(http.StatusOK, s.FailWarp("手机号码不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetSecureCodeData(string(tel))
	if err != nil {
		config.GetLogger().Warnw("发送手机验证码失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}

func (s *User) VerifyCodeHandler(c *gin.Context) {
	cont, _ := ioutil.ReadAll(c.Request.Body)

	if cont == nil {
		config.GetLogger().Warnw("信息不能为空",
			"cont:", cont,
		)
		c.JSON(http.StatusOK, s.FailWarp("信息不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetVerifyCodeData(string(cont))
	if err != nil {
		config.GetLogger().Warnw("验证手机验证码失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}

func (s *User) UserInfoHandler(c *gin.Context) {
	userID, _ := ioutil.ReadAll(c.Request.Body)

	if userID == nil {
		config.GetLogger().Warnw("用户ID不能为空",
			"userID:", userID,
		)
		c.JSON(http.StatusOK, s.FailWarp("用户ID不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetUserInfoData(string(userID))
	if err != nil {
		config.GetLogger().Warnw("获取用户基本信息失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}

func (s *User) ChangeInfoHandler(c *gin.Context) {
	cont, _ := ioutil.ReadAll(c.Request.Body)

	if cont == nil {
		config.GetLogger().Warnw("用户信息不能为空",
			"cont:", cont,
		)
		c.JSON(http.StatusOK, s.FailWarp("用户信息不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetChangeData(string(cont))
	if err != nil {
		config.GetLogger().Warnw("更新用户基本信息失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}

func (s *User) UserRunHandler(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		config.GetLogger().Warnw("用户ID不能为空",
			"userID:", userID,
		)
		c.JSON(http.StatusOK, s.FailWarp("用户ID不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetUserRunData(userID)
	if err != nil {
		config.GetLogger().Warnw("获取用户基本信息失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}

func (s *User) UserWalkHandler(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		config.GetLogger().Warnw("用户ID不能为空",
			"userID:", userID,
		)
		c.JSON(http.StatusOK, s.FailWarp("用户ID不能为空"))
		return
	}

	user := model.NewUser()

	err, data := user.GetUserRunData(userID)
	if err != nil {
		config.GetLogger().Warnw("获取用户基本信息失败",
			"err", err.Error(),
		)
		c.JSON(http.StatusOK, s.FailWarp(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.SuccessWarp(data))
}