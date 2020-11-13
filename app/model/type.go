package model

//用户实例
type User struct {
	Info         Info
	LoginData    LoginData
	RegisterData RegisterData
}

//用户基本信息
type Info struct {
	User          string `json:"user_name"`
	Password      string `json:"user_pwd"`
	Sex           string `json:"user_sex"`
	Tel           string `json:"user_tel"`
	Email         string `json:"user_email"`
	Type          string `json:"user_type"`
	Head_portrait string `json:"head_portrait"`
}

//用户登录返回字段
type LoginData struct {
	IsLogin bool `json:"is_login"`
	Type    int  `json:"type"`
}

//用户注册返回字段
type RegisterData struct {
	Registered bool `json:"registered"`
}
