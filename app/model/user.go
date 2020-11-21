package model

import (
	"back/app/config"
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
	}

	return temp
}

/**
 * 登录验证
 */
func (u *User) login(user_name string, user_pwd string) (err error) {
	data := &u.LoginData
	config.GetLogger().Info("开始登录")
	//查询用户类型
	row := db.Table("user_info").Where("user = ? AND password = ?", user_name, user_pwd).Select("type").Row()

	err = row.Scan(&data.Type)
	if err != nil {
		data.IsLogin = false
		config.GetLogger().Warnw("登录失败",
			"err", err.Error(),
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
func (u *User) register(i *Info) (err error) {
	data := &u.RegisterData
	count := 0
	config.GetLogger().Info("开始注册")
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
	//生成6位随机验证码
	config.GetLogger().Info("开始生成6位随机验证码")
	code := randCode()
	config.GetLogger().Info("生成6位随机验证码结束")

	config.GetLogger().Info("开始发送验证码")
	//检查用于发送验证码的手机号是否已经被注册
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4G4TXShUqRfEf1AnpaMx", "MH8TYZoKEJdnsgM63tSQQwMCIezKst")

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = tel
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
	config.GetRedis().Del(tel)
	config.GetRedis().Set(tel, code, 1*time.Minute)

	data.Sent = true

	return
}

/**
 * 验证码验证
 */
func (u *User) verify(tel string, code string) (err error) {
	data := &u.VerifyCodeData
	config.GetLogger().Info("开始获取redis中的验证码")

	tempCode, errs := config.GetRedis().Get(tel).Result()

	if errs != nil {
		data.Verified = false
		config.GetLogger().Warnw("验证手机验证码失败",
			"err", errs,
		)
		return
	}

	config.GetLogger().Info("获取redis中的验证码结束")
	println(tel, code)

	config.GetLogger().Info("开始校验验证码")
	if tempCode != code {
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
	config.GetLogger().Info("开始获取个人信息")

	row := db.Table("user_info").Where("id = ?", userID).
		Select("user, password, tel, email, sex, type, head_portrait").Row();
	
	err = row.Scan(&data.User, &data.Password, &data.Tel, &data.Email, &data.Sex, &data.Type, &data.Head_portrait)
	if err != nil {
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err,
		)
		return
	}

	config.GetLogger().Info("获取个人信息结束")

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

func (u *User) GetSecureCodeData(tel string) (err error, data SecureCodeData) {
	config.GetLogger().Info("开始发送手机验证码")

	err = u.send(tel)

	data = u.SecureCodeData

	config.GetLogger().Info("发送手机验证码结束")

	return
}

func (u *User) GetVerifyCodeData(tel string, code string) (err error, data VerifyCodeData) {
	config.GetLogger().Info("开始获取验证手机验证码数据")

	err = u.verify(tel, code)

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
