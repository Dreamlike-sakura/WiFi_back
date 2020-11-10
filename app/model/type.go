package model

import (
	"sync"
)

//用户实例
type User struct {
	wg        sync.WaitGroup
	Info      Info
	LoginData LoginData
}

//用户基本信息
type Info struct {
	Username     string `json:"user_name"`
	UserPassword string `json:"user_pwd"`
	UserSex      string `json:"user_sex"`
	UserTel      string `json:"user_tel"`
	UserEmail    string `json:"user_email"`
	UserType     string `json:"user_type"`
}

//用户登录返回字段
type LoginData struct {
	IsLogin bool `json:"is_login"`
	Type    int  `json:"type"`
}
