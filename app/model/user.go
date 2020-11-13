package model

import (
	"back/app/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/**
 * 构造函数, 得到实例
 */
func NewUser() *User {
	temp := &User{
		Info:      Info{},
		LoginData: LoginData{},
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

//----------------------------------分割线----------------------------------------
func (u *User) GetLoginData(user_name string, user_pwd string) (err error, data LoginData) {
	config.GetLogger().Info("开始获取登录数据")

	err = u.login(user_name, user_pwd)

	data = u.LoginData

	config.GetLogger().Info("获取登录数据结束")

	return
}
