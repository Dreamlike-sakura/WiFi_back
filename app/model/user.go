package model

import (
	"back/app/config"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"path"
	"strconv"
	"time"
)

/**
 * 构造函数, 得到实例
 */
func NewUser() *User {
	temp := &User{
		Info:              Info{},
		LoginData:         LoginData{},
		RegisterData:      RegisterData{},
		SecureCodeData:    SecureCodeData{},
		VerifyCodeData:    VerifyCodeData{},
		MovementData:      MovementData{},
		ModifyData:        ModifyData{},
		ChangePwdData:     ChangePwdData{},
		MovementListData:  []MovementListData{},
		CheckMovement:     CheckMovement{},
		CheckHeadPortrait: []CheckHeadPortrait{},
		GoPyData:          GoPyData{},
	}

	return temp
}

/**
 * 登录验证
 */
func (u *User) login(cont string) (err error) {
	data := &u.LoginData
	config.GetLogger().Info("开始登录")

	user := new(ReceiveLogin)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		data.IsLogin = false
		config.GetLogger().Warnw("登录数据解析失败",
			"err", err.Error,
		)
		return err
	}

	tempPwd := md5.Sum([]byte(user.UserPassword))
	md5str := fmt.Sprintf("%x", tempPwd)

	userName := user.UserName
	userPwd := md5str

	//查询用户类型
	row := db.Table("user_info").Where("user = ? AND password = ?", userName, userPwd).Select("type, id").Row()

	err = row.Scan(&data.Type, &data.UserID)
	if err != nil {
		data.IsLogin = false
		config.GetLogger().Warnw("登录失败",
			"err", err.Error,
		)
		return err
	}

	data.IsLogin = true
	config.GetLogger().Info("登录结束")

	return
}

/**
 * 用户注册
 */
func (u *User) register(cont string) (err error) {
	data := &u.RegisterData
	count := 0

	config.GetLogger().Info("开始解析注册数据")
	user := new(Info)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		data.Registered = false
		config.GetLogger().Warnw("数据解析失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("解析注册数据结束")

	config.GetLogger().Info("开始注册")
	//查询用户名是否重复，重复返回错误，否则数据库里插入 一条数据
	db.Table("user_info").Where("user = ?", user.User).Count(&count)

	//用户名存在时，
	if count != 0 {
		data.Registered = false
		config.GetLogger().Warnw("注册失败",
			"err", errors.New("用户名已存在"),
		)
		return errors.New("用户名已存在")
	}

	//密码进行MD5加密
	tempPwd := md5.Sum([]byte(user.Password))
	md5str := fmt.Sprintf("%x", tempPwd)

	user.Password = md5str
	user.Head_portrait = "1"
	user.Type = "0"
	//数据库中新建一个用户
	if err = db.Table("user_info").Create(user).Error; err != nil {
		data.Registered = false
		config.GetLogger().Warnw("注册失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("注册结束")

	data.Registered = true

	return
}

/**
 * 6位随机验证码
 */
func randCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rndCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	fmt.Println(rndCode)
	return rndCode
}

/**
 * 发送验证码
 */
func (u *User) send(tel string) (err error) {
	data := &u.SecureCodeData
	user := new(ReceiveTel)
	config.GetLogger().Info("开始解析数据")
	err = json.Unmarshal([]byte(tel), &user)
	if err != nil {
		config.GetLogger().Warnw("登录数据解析失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("解析数据结束")
	//生成6位随机验证码
	config.GetLogger().Info("开始生成6位随机验证码")
	code := randCode()
	config.GetLogger().Info("生成6位随机验证码结束")

	config.GetLogger().Info("开始发送验证码")
	//检查用于发送验证码的手机号是否已经被注册
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "", "")

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = user.UserTel
	request.SignName = "WIFI人体动作识别系统"
	request.TemplateCode = "SMS_205458618"
	request.TemplateParam = "{code:" + code + "}"
	//request.TemplateParam = "{code:123456}"

	response, err := client.SendSms(request)
	if err != nil {
		data.Sent = false
		config.GetLogger().Warnw("获取手机验证码失败",
			"err", err,
		)
		return
	}
	fmt.Printf("response is %#v\n", response)

	config.GetLogger().Info("发送验证码结束")
	//redis储存验证码，1分钟
	config.GetRedis().Del(user.UserTel)
	config.GetRedis().Set(user.UserTel, code, 1*time.Minute)

	data.Sent = true

	return
}

/**
 * 验证码验证
 */
func (u *User) verify(cont string) (err error) {
	data := &u.VerifyCodeData

	config.GetLogger().Info("开始解析数据")
	user := new(ReceiveTelAndCode)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		config.GetLogger().Warnw("登录数据解析失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("解析数据结束")

	config.GetLogger().Info("开始查询与此手机号绑定的用户ID")
	row := db.Table("user_info").Where("tel = ?", user.UserTel).Select("id").Row()
	err = row.Scan(&data.UserID)
	if err != nil {
		data.Verified = false
		config.GetLogger().Warnw("查询失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("查询与此手机号绑定的用户ID结束")

	config.GetLogger().Info("开始获取redis中的验证码")
	tempCode, errs := config.GetRedis().Get(user.UserTel).Result()
	if errs != nil {
		data.Verified = false
		config.GetLogger().Warnw("验证手机验证码失败",
			"err", errs,
		)
		return
	}
	config.GetLogger().Info("获取redis中的验证码结束")

	config.GetLogger().Info("开始校验验证码")
	if tempCode != user.UserSecureCode {
		data.Verified = false
		return errors.New("验证码错误")
	} else {
		data.Verified = true
	}
	config.GetLogger().Info("校验验证码结束")
	return
}

/**
 * 查看个人信息
 */
func (u *User) info(userID string) (err error) {
	data := &u.Info

	config.GetLogger().Info("开始解析数据")
	user := new(ReceiveID)
	err = json.Unmarshal([]byte(userID), &user)
	if err != nil {
		config.GetLogger().Warnw("数据解析失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("完成解析数据")

	config.GetLogger().Info("开始获取个人信息")

	row := db.Table("user_info, head_portrait").Where("id = ? AND head_portrait = picture_id", user.UserID).
		Select("id, user, password, tel, email, sex, type, url").Row()

	err = row.Scan(&data.ID, &data.User, &data.Password, &data.Tel, &data.Email, &data.Sex, &data.Type, &data.Head_portrait)
	if err != nil {
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err.Error,
		)
		return
	}

	config.GetLogger().Info("获取个人信息结束")

	return
}

/**
 * 修改个人信息
 */
func (u *User) changeInfo(cont string) (err error) {
	data := &u.ModifyData
	i := new(Info)
	count := 0

	config.GetLogger().Info("开始解析数据")
	user := new(ReceiveChange)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		data.Modified = false
		config.GetLogger().Warnw("数据解析失败",
			"err", err.Error(),
		)
		return err
	}
	config.GetLogger().Info("完成解析数据")

	config.GetLogger().Info("开始获取个人信息")
	err = db.Table("user_info").Where("id = ?", user.UserID).Count(&count).Error
	if err != nil || count == 0 {
		data.Modified = false
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("获取个人信息结束")

	i.User = user.UserName
	i.Sex = user.UserSex
	i.Tel = user.UserTel
	i.Email = user.UserEmail

	row := db.Raw(`SELECT picture_id FROM head_portrait WHERE url = ?`, user.HeadPortrait).Row()
	err = row.Scan(&i.Head_portrait)
	if err != nil {
		data.Modified = false
		config.GetLogger().Warnw("更新个人信息失败",
			"err", err,
		)
		return
	}

	err = db.Table("user_info").Model(&i).Where("id = ?", user.UserID).Updates(map[string]interface{}{"user": i.User, "sex": i.Sex, "tel": i.Tel, "email": i.Email, "head_portrait": i.Head_portrait}).Error
	if err != nil {
		data.Modified = false
		config.GetLogger().Warnw("更新个人信息失败",
			"err", err,
		)
		return
	}
	data.Modified = true

	config.GetLogger().Info("更新个人信息结束")

	return
}

/**
 * 修改个人密码
 */
func (u *User) changePwd(cont string) (err error) {
	data := &u.ChangePwdData
	i := new(Info)
	count := 0

	config.GetLogger().Info("开始解析数据")
	user := new(ReceiveChangePwd)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		data.Changed = false
		config.GetLogger().Warnw("数据解析失败",
			"err", err.Error(),
		)
		return err
	}
	config.GetLogger().Info("完成解析数据")

	config.GetLogger().Info("开始获取个人信息")
	err = db.Table("user_info").Where("id = ?", user.UserID).Count(&count).Error
	if err != nil || count == 0 {
		data.Changed = false
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("获取个人信息结束")

	config.GetLogger().Info("开始更新个人密码")

	tempPwd := md5.Sum([]byte(user.UserPassword))
	md5str := fmt.Sprintf("%x", tempPwd)

	i.Password = md5str

	err = db.Table("user_info").Model(&i).Where("id = ?", user.UserID).Updates(map[string]interface{}{"password": i.Password}).Error
	if err != nil {
		data.Changed = false
		config.GetLogger().Warnw("更新个人密码失败",
			"err", err,
		)
		return
	}
	data.Changed = true

	config.GetLogger().Info("更新个人密码结束")

	return
}

func (u *User) changePwd2(cont string) (err error) {
	data := &u.ChangePwdData
	i := new(Info)
	count := 0
	temp := ""

	config.GetLogger().Info("开始解析数据")
	user := new(ReceiveChangePwd2)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		data.Changed = false
		config.GetLogger().Warnw("数据解析失败",
			"err", err.Error(),
		)
		return err
	}
	config.GetLogger().Info("完成解析数据")

	config.GetLogger().Info("开始获取个人信息")
	tempOldPwd := md5.Sum([]byte(user.OldPassword))
	md5str := fmt.Sprintf("%x", tempOldPwd)
	err = db.Table("user_info").Where("id = ? AND password = ?", user.ID, md5str).Count(&count).Error
	if err != nil || count == 0 {
		data.Changed = false
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("获取个人信息结束")

	config.GetLogger().Info("开始检验密码是否重复")
	tempNewPwd := md5.Sum([]byte(user.NewPassword))
	md5str = fmt.Sprintf("%x", tempNewPwd)

	row := db.Table("user_info").Where("id = ?", user.ID).Select("password").Row()
	err = row.Scan(&temp)

	if err != nil {
		data.Changed = false
		config.GetLogger().Warnw("检验密码是否重复失败",
			"err", err,
		)
		return
	}
	if temp == md5str {
		data.Changed = false
		config.GetLogger().Warnw("新密码与原密码重复")
		return errors.New("新密码与原密码重复")
	}
	config.GetLogger().Info("检验密码是否重复结束")

	config.GetLogger().Info("开始更新个人密码")
	i.Password = md5str

	err = db.Table("user_info").Model(&i).Where("id = ?", user.ID).Updates(map[string]interface{}{"password": i.Password}).Error
	if err != nil {
		data.Changed = false
		config.GetLogger().Warnw("更新个人密码失败",
			"err", err,
		)
		return
	}
	data.Changed = true
	config.GetLogger().Info("更新个人密码结束")

	return
}

/**
 * 查看用户动作信息列表
 */
func (u *User) movementList(cont string) (err error) {
	data := &u.MovementListData

	config.GetLogger().Info("开始解析注册数据")
	user := new(ReceiveMovementList)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		config.GetLogger().Warnw("数据解析失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("解析注册数据结束")

	config.GetLogger().Info("开始获取动作数据并分页")

	table := [3]string{"dealt_run", "dealt_walk", "dealt_shakehand"}
	tempType, err := strconv.Atoi(user.Type)
	if err != nil {
		config.GetLogger().Warnw("类型错误",
			"err:", err,
		)
		return
	}
	if tempType < 1 || tempType > 3 {
		config.GetLogger().Warnw("类型错误",
			"err:", errors.New("类型范围错误"),
		)
		return
	}

	rows, err := db.Table(table[tempType-1]).Where("uid = ?", user.UserID).
		Order("time").Limit(user.PageNum * user.PageSize).Offset((user.PageNum - 1) * user.PageSize).Select("id, time").Rows()
	if err != nil {
		config.GetLogger().Warnw("数据库数据错误",
			"err:", err,
		)
		return err
	}

	defer rows.Close()
	for rows.Next() {
		tempData := new(MovementListData)
		tempData.ID = ""
		tempData.Type = user.Type
		tempData.FileName = ""
		tempData.Time = ""

		err = rows.Scan(&tempData.ID, &tempData.Time)
		if err != nil {
			config.GetLogger().Warnw("赋值错误",
				"err:", err,
			)
			return err
		}

		tempData.FileName = tempData.Time

		*data = append(*data, *tempData)
	}

	config.GetLogger().Info("获取动作数据并分页结束")

	return
}

/**
 * 查看头像列表信息
 */
func (u *User) headPortraitList() (err error) {
	data := &u.CheckHeadPortrait

	config.GetLogger().Info("开始获取头像信息")

	rows, errs := db.Table("head_portrait").Select("picture_id, url").Rows()
	if errs != nil {
		config.GetLogger().Warnw("获取头像信息失败",
			"err", errs.Error,
		)
		return
	}

	for rows.Next() {
		temp := new(CheckHeadPortrait)
		temp.ID = ""
		temp.Url = ""

		err = rows.Scan(&temp.ID, &temp.Url)
		if err != nil {
			config.GetLogger().Warnw("数据获取错误",
				"err:", err,
			)
			return err
		}

		*data = append(*data, *temp)
	}

	config.GetLogger().Info("获取头像信息结束")

	return
}

/**
 * go调用python
 */
func (u *User) goPy(cont string) (err error) {
	data := &u.GoPyData
	config.GetLogger().Info("开始解析原始数据")
	user := new(ReceiveGoPyData)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		config.GetLogger().Warnw("数据解析失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("解析原始数据结束")

	args := []string{"read_bfee_file.py", user.File}

	fmt.Println(args)

	cmd := exec.Command("python", args...)
	cmd.Dir = ".\\py\\"

	err = cmd.Run()
	if err != nil {
		data.Success = false
		config.GetLogger().Warnw("执行python脚本失败",
			"err", err,
		)
		return
	}
	data.Success = true

	return
}

/**
 * go读取文件
 */
func (u *User) getAmpOrPhase(cont string) (err error) {
	data := &u.CheckMovement

	config.GetLogger().Info("开始解析文件名")
	user := new(ReceiveCheckMovement)
	err = json.Unmarshal([]byte(cont), &user)
	if err != nil {
		config.GetLogger().Warnw("文件名解析失败",
			"err", err.Error,
		)
		return err
	}
	config.GetLogger().Info("解析文件名结束")

	//dir := "D:\20study\2020project\back\"
	dir := ".\\data\\wifi\\"
	fileStr := path.Join(dir, user.FileName)
	//fileStr := dir + user.FileName

	fmt.Println(fileStr)

	f, errs := ioutil.ReadFile(fileStr)
	if errs != nil {
		fmt.Println("read fail", errs)
	}

	config.GetLogger().Info("开始解析矩阵数据")
	err = json.Unmarshal([]byte(f), &data.Content)
	if err != nil {
		config.GetLogger().Warnw("矩阵数据解析失败",
			"err", err,
		)
		return err
	}
	config.GetLogger().Info("矩阵数据解析结束")

	return
}

/**
 * 上传文件
 */
//func (u *User) upload (c *gin.Context) (err error) {
//	config.GetLogger().Info("开始获取文件")
//	file, header, err := c.Request.FormFile("file")
//
//	if err != nil {
//		config.GetLogger().Warnw("文件读取失败",
//			"err", err.Error,
//		)
//		return err
//	}
//	config.GetLogger().Info("获取文件结束")
//
//	return
//}

//----------------------------------分割线----------------------------------------
func (u *User) GetLoginData(cont string) (err error, data LoginData) {
	config.GetLogger().Info("开始获取登录数据")

	err = u.login(cont)

	data = u.LoginData

	config.GetLogger().Info("获取登录数据结束")

	return
}

func (u *User) GetRegisterData(cont string) (err error, data RegisterData) {
	config.GetLogger().Info("开始获取注册数据")

	err = u.register(cont)

	data = u.RegisterData

	config.GetLogger().Info("获取注册数据结束")

	return
}

func (u *User) GetSecureCodeData(tel string) (err error, data SecureCodeData) {
	config.GetLogger().Info("开始发送手机验证码")

	err = u.send(tel)

	data = u.SecureCodeData

	config.GetLogger().Info("发送手机验证码结束")

	return
}

func (u *User) GetVerifyCodeData(cont string) (err error, data VerifyCodeData) {
	config.GetLogger().Info("开始获取验证手机验证码数据")

	err = u.verify(cont)

	data = u.VerifyCodeData

	config.GetLogger().Info("获取验证手机验证码数据结束")

	return
}

func (u *User) GetUserInfoData(userID string) (err error, data Info) {
	config.GetLogger().Info("开始获取用户基本信息数据")

	err = u.info(userID)

	data = u.Info

	config.GetLogger().Info("获取用户基本信息数据结束")

	return
}

func (u *User) GetChangeData(cont string) (err error, data ModifyData) {
	config.GetLogger().Info("开始修改用户信息")

	err = u.changeInfo(cont)

	data = u.ModifyData

	config.GetLogger().Info("修改用户信息结束")

	return
}

func (u *User) GetChangePwdData(cont string) (err error, data ChangePwdData) {
	config.GetLogger().Info("开始修改用户密码")

	err = u.changePwd(cont)

	data = u.ChangePwdData

	config.GetLogger().Info("修改用户信息密码结束")

	return
}

func (u *User) GetChangePwdData2(cont string) (err error, data ChangePwdData) {
	config.GetLogger().Info("开始修改用户密码")

	err = u.changePwd2(cont)

	data = u.ChangePwdData

	config.GetLogger().Info("修改用户信息密码结束")

	return
}

func (u *User) GetMovementListData(cont string) (err error, data []MovementListData) {
	config.GetLogger().Info("开始查询用户动作列表")

	err = u.movementList(cont)

	data = u.MovementListData

	config.GetLogger().Info("查询用户动作列表结束")

	return
}

func (u *User) GetAPData(cont string) (err error, data CheckMovement) {
	config.GetLogger().Info("开始查询用户动作")

	err = u.getAmpOrPhase(cont)

	data = u.CheckMovement

	config.GetLogger().Info("查询用户动作列表")

	return
}

func (u *User) GetHeadPortraitData() (err error, data []CheckHeadPortrait) {
	config.GetLogger().Info("开始查询头像列表")

	err = u.headPortraitList()

	data = u.CheckHeadPortrait

	config.GetLogger().Info("查询头像列表结束")

	return
}

func (u *User) GetGoPyData(cont string) (err error, data GoPyData) {
	config.GetLogger().Info("开始调用python")

	err = u.goPy(cont)

	data = u.GoPyData

	config.GetLogger().Info("调用python结束")

	return
}

