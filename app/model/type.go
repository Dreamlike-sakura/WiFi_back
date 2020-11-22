package model

//用户实例
type User struct {
	Info           Info
	LoginData      LoginData
	RegisterData   RegisterData
	SecureCodeData SecureCodeData
	VerifyCodeData VerifyCodeData
	MovementData   MovementData
	ModifyData     ModifyData
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

//用于接受用户登录时的参数
type ReceiveLogin struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_pwd"`
}

//用户注册返回字段
type RegisterData struct {
	Registered bool `json:"registered"`
}

//手机验证码返回字段
type SecureCodeData struct {
	Sent bool `json:"sent"`
}

//确认验证码返回字段
type VerifyCodeData struct {
	Verified bool `json:"identify"`
}

//修改个人信息返回字段
type ModifyData struct {
	Modified bool `json:"modified"`
}

//动作信息返回
type MovementData struct {
	DealtAmplitude string `json:"dealt_amplitude"`
	DealtPhase     string `json:"dealt_phase"`
	Amplitude      string `json:"amplitude"`
	Phase          string `json:"phase"`
	Abnormal       string `json:"abnormal"`
	Time           string `json:"time"`
}
