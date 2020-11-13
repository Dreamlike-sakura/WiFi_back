package model

import (
	"back/app/config"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/**
 * 构造函数, 得到实例
 */
func NewUser() *User {
	temp := &User{
		Info:         Info{},
		LoginData:    LoginData{},
		RegisterData: RegisterData{},
	}

	return temp
}

/**
 * 登录验证
 */
func (u *User) login(user_name string, user_pwd string) (err error) {
	data := &u.LoginData
	//查询用户类型
	row := db.Table("user_info").Where("user = ? AND pwd = ?", user_name, user_pwd).Select("type").Row()

	err = row.Scan(&data.Type)
	if err != nil {
		data.IsLogin = false
		config.GetLogger().Warnw("登录失败",
			"err", err.Error(),
		)
		return err
	}

	data.IsLogin = true

	return
}

/**
 * 用户注册
 */
func (u *User) register(i *Info) (err error) {
	data := &u.RegisterData
	count := 0
	//查询用户名是否重复，重复返回错误，否则数据库里插入 一条数据
	db.Table("user_info").Where("user = ?", i.User).Count(&count)

	//用户名存在时，
	if count != 0 {
		data.Registered = false
		config.GetLogger().Warnw("注册失败",
			"err", errors.New("用户名已存在"),
		)
		return errors.New("用户名已存在")
	}

	//数据库中新建一个用户
	if err = db.Table("user_info").Create(i).Error; err != nil {
		data.Registered = false
		config.GetLogger().Warnw("注册失败",
			"err", errors.New("新建用户失败"),
		)
		return
	}

	data.Registered = true

	return
}

//----------------------------------分割线----------------------------------------
func (u *User) GetLoginData(user_name string, user_pwd string) (err error, data LoginData) {
	config.GetLogger().Info("开始获取登录数据")

	err = u.login(user_name, user_pwd)

	data = u.LoginData

	config.GetLogger().Info("获取登录数据结束")

	return
}

func (u *User) GetRegisterData(i *Info) (err error, data RegisterData) {
	config.GetLogger().Info("开始获取注册数据")

	err = u.register(i)

	data = u.RegisterData

	config.GetLogger().Info("获取注册数据结束")

	return
}
