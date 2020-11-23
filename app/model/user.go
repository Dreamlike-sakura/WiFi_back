package model

import (
	"back/app/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"math/rand"
	"time"
)

/**
 * 构造函数, 得到实例
 */
func NewUser() *User {
	temp := &User{
		Info:           Info{},
		LoginData:      LoginData{},
		RegisterData:   RegisterData{},
		SecureCodeData: SecureCodeData{},
		VerifyCodeData: VerifyCodeData{},
		MovementData:   MovementData{},
		ModifyData:     ModifyData{},
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
	userName := user.UserName
	userPwd := user.UserPassword

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
			"err", err.Error(),
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

	user.Head_portrait = "1"
	user.Type = "0"
	//数据库中新建一个用户
	if err = db.Table("user_info").Create(user).Error; err != nil {
		data.Registered = false
		config.GetLogger().Warnw("注册失败",
			"err", errors.New("新建用户失败"),
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
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4G4TXShUqRfEf1AnpaMx", "MH8TYZoKEJdnsgM63tSQQwMCIezKst")

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
			"err", err.Error(),
		)
		return err
	}
	config.GetLogger().Info("完成解析数据")

	config.GetLogger().Info("开始获取个人信息")

	row := db.Table("user_info").Where("id = ?", user.UserID).
		Select("id, user, password, tel, email, sex, type, head_portrait").Row()

	err = row.Scan(&data.ID, &data.User, &data.Password, &data.Tel, &data.Email, &data.Sex, &data.Type, &data.Head_portrait)
	if err != nil {
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err,
		)
		return
	}

	config.GetLogger().Info("获取个人信息结束")

	return
}

/**
 * 查看个人动作信息之跑步
 */
func (u *User) runInfo(userID string) (err error) {
	data := &u.MovementData
	config.GetLogger().Info("开始查询跑步动作信息")

	row := db.Table("dealt_run").Where("uid = ?", userID).
		Select("origin_amplitude, amplitude, origin_phase, phase, abnormal, time").Row()

	err = row.Scan(&data.Amplitude, &data.DealtAmplitude, &data.Phase, &data.DealtPhase, &data.Abnormal, &data.Time)
	if err != nil {
		config.GetLogger().Warnw("查询跑步动作信息失败",
			"err", err,
		)
		return
	}

	config.GetLogger().Info("查询跑步动作信息结束")

	return
}

/**
 * 查看个人动作信息之行走
 */
func (u *User) walkInfo(userID string) (err error) {
	data := &u.MovementData
	config.GetLogger().Info("开始查询行走动作信息")

	row := db.Table("dealt_walk").Where("uid = ?", userID).
		Select("origin_amplitude, amplitude, origin_phase, phase, abnormal, time").Row()

	err = row.Scan(&data.Amplitude, &data.DealtAmplitude, &data.Phase, &data.DealtPhase, &data.Abnormal, &data.Time)
	if err != nil {
		config.GetLogger().Warnw("查询行走动作信息失败",
			"err", err,
		)
		return
	}

	config.GetLogger().Info("查询行走动作信息结束")

	return
}

/**
 * 查看个人动作信息之摇手
 */
func (u *User) shakeInfo(userID string) (err error) {
	data := &u.MovementData
	config.GetLogger().Info("开始查询摇手动作信息")

	row := db.Table("dealt_shakehand").Where("uid = ?", userID).
		Select("origin_amplitude, amplitude, origin_phase, phase, abnormal, time").Row()

	err = row.Scan(&data.Amplitude, &data.DealtAmplitude, &data.Phase, &data.DealtPhase, &data.Abnormal, &data.Time)
	if err != nil {
		config.GetLogger().Warnw("查询摇手动作信息失败",
			"err", err,
		)
		return
	}

	config.GetLogger().Info("查询摇手动作信息结束")

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
	err = db.Table("user_info").Where("id = ?", user.UserID).Count(&count).Error;
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
	i.Head_portrait = user.HeadPortrait

	err = db.Table("user_info").Model(&i).Where("id = ?", user.UserID).Updates(map[string]interface{}{"user":i.User, "sex":i.Sex, "tel":i.Tel, "email":i.Email, "head_portrait":i.Head_portrait}).Error
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

//----------------------------------分割线----------------------------------------
func (u *User) GetLoginData(cont string) (err error, data LoginData) {
	config.GetLogger().Info("开始获取登录数据")

	err = u.login(cont)

	data = u.LoginData

	config.GetLogger().Info("获取登录数据结束")

	return
}

func (u *User) GetRegisterData( cont string) (err error, data RegisterData) {
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

func (u *User) GetUserRunData(userID string) (err error, data MovementData) {
	config.GetLogger().Info("开始获取用户跑步信息数据")

	err = u.runInfo(userID)

	data = u.MovementData

	config.GetLogger().Info("获取用户跑步信息数据结束")

	return
}

func (u *User) GetUserWalkData(userID string) (err error, data MovementData) {
	config.GetLogger().Info("开始获取用户跑步信息数据")

	err = u.walkInfo(userID)

	data = u.MovementData

	config.GetLogger().Info("获取用户跑步信息数据结束")

	return
}

func (u *User) GetUserShakeData(userID string) (err error, data MovementData) {
	config.GetLogger().Info("开始获取用户摇手信息数据")

	err = u.shakeInfo(userID)

	data = u.MovementData

	config.GetLogger().Info("获取用户摇手信息数据结束")

	return
}

func (u *User) GetChangeData(cont string) (err error, data ModifyData) {
	config.GetLogger().Info("开始修改用户信息")

	err = u.changeInfo(cont)

	data = u.ModifyData

	config.GetLogger().Info("修改用户信息结束")

	return
}
