package model

//用户实例
type User struct {
	Info              Info
	LoginData         LoginData
	RegisterData      RegisterData
	SecureCodeData    SecureCodeData
	VerifyCodeData    VerifyCodeData
	MovementData      MovementData
	ModifyData        ModifyData
	ChangePwdData     ChangePwdData
	MovementListData  MovementListData
	CheckMovement     CheckMovement
	CheckHeadPortrait []CheckHeadPortrait
	GoPyData          GoPyData
	StatisticsData    []StatisticsData
	UploadData        UploadData
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

//用于接受用户修改密码时的参数
type ReceiveChangePwd2 struct {
	ID          string `json:"id"`
	OldPassword string `json:"old_pwd"`
	NewPassword string `json:"new_pwd"`
}

//用于接受用户ID
type ReceiveID struct {
	UserID string `json:"user_id"`
}

//用于读取ak
type ReceiveAK struct {
	AK  string `json:"ak"`
	AKS string `json:"aks"`
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
	UserName string `json:"user_name"`
	FileName string `json:"file_name"`
	FileType int    `json:"file_type"` //1是跑步、2是行走、3是搖手
	MoveType string `json:"move_type"` //要读取文件的类型“amp”幅度、“phase”相位、“abnormal”异常
	Type     int    `json:"type"`      //是否是去噪后的，0为原始数据，1为去噪后数据
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

//用户上传文件返回字段
type UploadData struct {
	Uploaded bool `json:"uploaded"`
}

//接收用户调用python的参数
type ReceiveGoPyData struct {
	ID   string `json:"id"`
	File string `json:"file"`
}

//手机验证码返回字段
type SecureCodeData struct {
	Sent bool `json:"sent"`
}

//调用Python返回字段
type GoPyData struct {
	Success bool `json:"success"`
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

//用户查看头像列表返回字段
type CheckHeadPortrait struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

//统计数据
type StatisticsData struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
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
type MoveData struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	FileName string `json:"file_name"`
	Time     string `json:"time"`
}

//用户动作信息返回列表
type MovementListData struct {
	Sum  int        `json:"sum"`
	List []MoveData `json:"list"`
}

//往数据库里插入原始数据
type OriginData struct {
	UID      string `gorm:"Column:uid;Type:int;NOT NULL"`
	FileName string `gorm:"Column:filename;Type:varchar(255);DEFAULT:NULL"`
	DataUrl  string `gorm:"Column:data_url;Type:varchar(255);DEFAULT:NULL"`
	Time     string `gorm:"Column:time;Type:datetime;NOT NULL"`
}

//往数据库里插入处理后数据
type DealtData struct {
	UID         string `gorm:"Column:uid;Type:int;NOT NULL"`
	FileName    string `gorm:"Column:filename;Type:varchar(255);NOT NULL"`
	Amp         string `gorm:"Column:amplitude;Type:text;DEFAULT:NULL"`
	Phase       string `gorm:"Column:phase;Type:text;DEFAULT:NULL"`
	OriginAmp   string `gorm:"Column:origin_amplitude;Type:text;DEFAULT:NULL"`
	OriginPhase string `gorm:"Column:origin_phase;Type:text;DEFAULT:NULL"`
	Abnormal    string `gorm:"Column:abnormal;Type:text;DEFAULT:NULL"`
	Time        string `gorm:"Column:time;Type:datetime;NOT NULL"`
}
