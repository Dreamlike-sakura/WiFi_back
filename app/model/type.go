package model

//用户实例
type User struct {
	Info             Info
	LoginData        LoginData
	RegisterData     RegisterData
	SecureCodeData   SecureCodeData
	VerifyCodeData   VerifyCodeData
	MovementData     MovementData
	ModifyData       ModifyData
	ChangePwdData    ChangePwdData
	MovementListData []MovementListData
	CheckMovement    CheckMovement
}

//用户基本信息
type Info struct {
	ID            string `json:"user_id"`
	User          string `json:"user_name"`
	Password      string `json:"user_pwd"`
	Sex           string `json:"user_sex"`
	Tel           string `json:"user_tel"`
	Email         string `json:"user_email"`
	Type          string `json:"user_type"`
	Head_portrait string `json:"head_portrait"`
}

//用于接受用户登录时的参数
type ReceiveLogin struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_pwd"`
}

//用于接受用户修改密码时的参数
type ReceiveChangePwd struct {
	UserID       string `json:"user_id"`
	UserPassword string `json:"user_pwd"`
}

//用于接受用户ID
type ReceiveID struct {
	UserID string `json:"user_id"`
}

//用于接受用户电话
type ReceiveTel struct {
	UserTel string `json:"tel"`
}

//用于接受用户验证码验证
type ReceiveTelAndCode struct {
	UserTel        string `json:"tel"`
	UserSecureCode string `json:"security_code"`
}

//用于接受用户查看动作信息时的参数
type ReceiveCheckMovement struct {
	FileName string `json:"file_name"`
}

//用于接受用户查看动作列表
type ReceiveMovementList struct {
	UserID   string `json:"user_id"`
	Type     string `json:"type"`
	PageNum  int    `json:"page_num"`
	PageSize int    `json:"page_size"`
}

//用于接受用户注册时的参数
type ReceiveRegister struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_pwd"`
	UserTel      string `json:"user_tel"`
	UserEmail    string `json:"user_email"`
}

//用于接受用户修改信息时的参数
type ReceiveChange struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	UserTel      string `json:"user_tel"`
	UserEmail    string `json:"user_email"`
	UserSex      string `json:"user_sex"`
	HeadPortrait string `json:"head_portrait"`
}

//用户登录返回字段
type LoginData struct {
	IsLogin bool   `json:"is_login"`
	Type    int    `json:"type"`
	UserID  string `json:"user_id"`
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
	UserID   string `json:"user_id"`
	Verified bool   `json:"identify"`
}

//修改个人信息返回字段
type ModifyData struct {
	Modified bool `json:"modified"`
}

//修改个人密码返回字段
type ChangePwdData struct {
	Changed bool `json:"changed"`
}

//用户查看动作信息返回字段
type CheckMovement struct {
	Content [][]float64 `json:"content"`
}

//动作信息返回字段
type MovementData struct {
	DealtAmplitude string `json:"dealt_amplitude"`
	DealtPhase     string `json:"dealt_phase"`
	Amplitude      string `json:"amplitude"`
	Phase          string `json:"phase"`
	Abnormal       string `json:"abnormal"`
	Time           string `json:"time"`
}

//用户动作信息返回列表
type MovementListData struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	FileName string `json:"file_name"`
	Time     string `json:"time"`
}
