package model

import (
	"back/app/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

/**
 * 构造函数, 得到实例
 */
func NewUser(username string) *User {
	temp := &User{
		wg:         sync.WaitGroup{},
		Info:       Info{},
		LoginData:  LoginData{},
	}

	temp.wg.Add(1)
	
	return temp
}

/**
 * 登录验证
 */
func (u *User) login(user_name string, user_pwd string) (err error) {
	data := &u.LoginData
	//查询用户类型
	//row := db.Table("user_info").Where("user = ?", user_name, "pwd = ?", user_pwd).Select("type").Row()
	row := db.Raw(
		`
			SELECT
				type
			FROM
				user_info
			WHERE
				user = ?
			AND
				pwd = ?`,
		user_name,
		user_pwd,
		).Row()
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
	go func() {
		temp := u.login(user_name, user_pwd)
		if temp != nil {
			err = temp
		}
		u.wg.Done()
	}()

	// 等待执行完成
	u.wg.Wait()
	data = u.LoginData

	config.GetLogger().Info("获取登录数据结束")

	return
}