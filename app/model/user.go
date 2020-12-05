package model

import (
	"back/app/config"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
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
		MovementListData:  MovementListData{},
		CheckMovement:     CheckMovement{},
		CheckHeadPortrait: []CheckHeadPortrait{},
		GoPyData:          GoPyData{},
		StatisticsData:    []StatisticsData{},
		UploadData:        UploadData{},
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

	//密码加密后验证
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
	//查询用户名是否重复，重复返回错误，否则数据库里插入一条数据
	db.Table("user_info").Where("user = ?", user.User).Count(&count)

	//用户名存在时，
	if count != 0 {
		data.Registered = false
		config.GetLogger().Warnw("注册失败",
			"err", errors.New("用户名已存在"),
		)
		return errors.New("用户名已存在")
	}

	//电话号码不能重复
	count = 0
	db.Table("user_info").Where("tel = ?", user.Tel).Count(&count)

	//号码存在时
	if count != 0 {
		data.Registered = false
		config.GetLogger().Warnw("注册失败",
			"err", errors.New("号码已存在"),
		)
		return errors.New("号码已存在")
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

	config.GetLogger().Info("开始创建用户目录")
	err = createDir(user.User)
	if err != nil {
		config.GetLogger().Warnw("文件夹创建失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("创建用户目录结束")

	data.Registered = true

	config.GetLogger().Info("注册结束")

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

	//数据解析
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

	config.GetLogger().Info("开始读取ak")
	fileStr := ".\\ak\\ak.json"

	f, err := ioutil.ReadFile(fileStr)
	if err != nil {
		fmt.Println("read fail", err)
		config.GetLogger().Warnw("文件读取失败",
			"err", err,
		)
		return
	}

	i := new(ReceiveAK)
	err = json.Unmarshal([]byte(f), &i)
	if err != nil {
		config.GetLogger().Warnw("矩阵数据解析失败",
			"err", err,
		)
		return err
	}
	config.GetLogger().Info("读取ak结束")

	config.GetLogger().Info("开始发送验证码")
	//检查用于发送验证码的手机号是否已经被注册
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", i.AK, i.AKS)
	if err != nil {
		data.Sent = false
		config.GetLogger().Warnw("获取手机验证码失败",
			"err", err,
		)
		return
	}

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

	//接收前端所传数据并解析
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

	//查询绑定的手机号
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

	//获取发送验证码时储存在redis中的验证码
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

	//校验验证码
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

	//接收前端所传数据并解析
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

	//查询个人信息
	config.GetLogger().Info("开始获取个人信息")
	row := db.Table("user_info, head_portrait").Where("id = ? AND head_portrait = picture_id", user.UserID).
		Select("id, user, password, tel, email, sex, type, url").Row()

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
 * 修改个人信息
 */
func (u *User) changeInfo(cont string) (err error) {
	data := &u.ModifyData
	i := new(Info)
	count := 0

	//接收前端所传数据并解析
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

	//查询用户是否存在
	config.GetLogger().Info("开始获取个人信息")
	err = db.Table("user_info").Where("id = ?", user.UserID).Count(&count).Error
	if err != nil || count == 0 {
		data.Modified = false
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err,
		)
		return errors.New("获取个人信息失败")
	}
	config.GetLogger().Info("获取个人信息结束")

	//用户名不可重复
	config.GetLogger().Info("开始检验用户名是否重复")
	count = 0
	db.Table("user_info").Where("user = ? AND id <> ?", user.UserName, user.UserID).Count(&count)

	//用户名存在时
	if count != 0 {
		data.Modified = false
		config.GetLogger().Warnw("更改信息失败",
			"err", errors.New("用户名已存在"),
		)
		return errors.New("用户名已存在")
	}
	config.GetLogger().Info("检验用户名是否重复结束")

	//电话号码不可重复
	config.GetLogger().Info("开始检验电话是否重复")
	count = 0
	db.Table("user_info").Where("tel = ? AND id <> ?", user.UserTel, user.UserID).Count(&count)

	//用户名存在时，
	if count != 0 {
		data.Modified = false
		config.GetLogger().Warnw("更改信息失败",
			"err", errors.New("电话号码已存在"),
		)
		return errors.New("电话号码已存在")
	}
	config.GetLogger().Info("检验号码是否重复结束")

	i.User = user.UserName
	i.Sex = user.UserSex
	i.Tel = user.UserTel
	i.Email = user.UserEmail

	config.GetLogger().Info("开始查询头像id")
	//根据前端发送的url找到对应的头像id，准备插入
	row := db.Raw(`SELECT picture_id FROM head_portrait WHERE url = ?`, user.HeadPortrait).Row()
	err = row.Scan(&i.Head_portrait)
	if err != nil {
		data.Modified = false
		config.GetLogger().Warnw("更新个人信息失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("查询头像id结束")

	config.GetLogger().Info("开始更新数据库数据")
	//开始更新
	err = db.Table("user_info").Model(&i).Where("id = ?", user.UserID).Updates(map[string]interface{}{"user": i.User, "sex": i.Sex, "tel": i.Tel, "email": i.Email, "head_portrait": i.Head_portrait}).Error
	if err != nil {
		data.Modified = false
		config.GetLogger().Warnw("更新个人信息失败",
			"err", err,
		)
		return
	}
	data.Modified = true
	config.GetLogger().Info("更新数据库数据结束")

	config.GetLogger().Info("更新个人信息结束")

	return
}

/**
 * 修改个人密码（忘记密码）
 */
func (u *User) changePwd(cont string) (err error) {
	data := &u.ChangePwdData
	i := new(Info)
	count := 0

	//接收前端所传数据并解析
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

	//检测用户是否存在
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

	//更新密码
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

/**
 * 修改个人密码（正常修改密码）
 */
func (u *User) changePwd2(cont string) (err error) {
	data := &u.ChangePwdData
	i := new(Info)
	count := 0
	temp := ""

	//接收前端所传数据并解析
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

	//获取用户信息
	config.GetLogger().Info("开始获取个人信息")
	tempOldPwd := md5.Sum([]byte(user.OldPassword))
	md5str := fmt.Sprintf("%x", tempOldPwd)

	err = db.Table("user_info").Where("id = ? AND password = ?", user.ID, md5str).Count(&count).Error
	if err != nil || count == 0 {
		data.Changed = false
		config.GetLogger().Warnw("获取个人信息失败",
			"err", err,
		)
		return errors.New("获取个人信息失败")
	}
	config.GetLogger().Info("获取个人信息结束")

	//新密码不可与原密码重复
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

	//更新密码
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

	//接收前端所传数据并解析
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

	//1--3分别表示跑步、行走、摇手
	if tempType < 1 || tempType > 3 {
		config.GetLogger().Warnw("类型错误",
			"err:", errors.New("类型范围错误"),
		)
		return
	}

	//获取动作数据，并根据页大小和页码进行分页
	rows, err := db.Table(table[tempType-1]).Where("uid = ?", user.UserID).
		Order("time").Limit(user.PageSize).Offset((user.PageNum - 1) * user.PageSize).Select("id, filename, time").Rows()
	if err != nil {
		config.GetLogger().Warnw("数据库数据错误",
			"err:", err,
		)
		return err
	}

	//查询数据总条数
	err = db.Table(table[tempType-1]).Where("uid = ?", user.UserID).Count(&data.Sum).Error
	if err != nil {
		config.GetLogger().Warnw("数据库数据错误",
			"err:", err,
		)
		return err
	}

	defer rows.Close()
	for rows.Next() {
		tempData := new(MoveData)
		tempData.ID = ""
		tempData.Type = user.Type
		tempData.FileName = ""
		tempData.Time = ""

		err = rows.Scan(&tempData.ID, &tempData.FileName, &tempData.Time)
		if err != nil {
			config.GetLogger().Warnw("赋值错误",
				"err:", err,
			)
			return err
		}

		data.List = append(data.List, *tempData)
	}

	config.GetLogger().Info("获取动作数据并分页结束")

	return
}

/**
 * 查看头像列表信息
 */
func (u *User) headPortraitList() (err error) {
	data := &u.CheckHeadPortrait

	//查询头像列表
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
	name := ""

	//接收前端所传数据并解析
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

	config.GetLogger().Info("开始查询用户名")
	row := db.Table("user_info").Where("id = ?", user.ID).Select("user").Row()
	err = row.Scan(&name)
	if err != nil {
		config.GetLogger().Warnw("用户名获取失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("查询用户名结束")

	//去掉文件后缀名
	fileSuffix := path.Ext(user.File)
	user.File = strings.TrimSuffix(user.File, fileSuffix)

	args := []string{"read_bfee_file.py", user.File, name}

	fmt.Println(args)

	//使用命令行的方式运行python文件
	cmd := exec.Command("python", args...)
	//设置执行文件的文件路径
	cmd.Dir = ".\\py\\"
	//开始执行
	err = cmd.Run()
	if err != nil {
		data.Success = false
		config.GetLogger().Warnw("执行python脚本失败",
			"err", err,
		)
		return
	}
	data.Success = true
	//脚本执行结束后会生成2个json文件

	return
}

/**
 * go读取文件
 */
func (u *User) getAmpOrPhase(cont string) (err error) {
	data := &u.CheckMovement
	key1 := [3]string{"run", "walk", "shake"}
	key2 := [2]string{"origin", "dealt"}

	//接收前端所传数据并解析
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

	config.GetLogger().Info("开始进行数据有效性验证")
	if user.MoveType != "amp" && user.MoveType != "phase" && user.MoveType != "abnormal" {
		config.GetLogger().Warnw("MoveType错误",
			"err", errors.New("MoveType错误"),
		)
		return errors.New("MoveType错误")
	}

	if user.FileType < 1 || user.FileType > 3 {
		config.GetLogger().Warnw("FileType错误",
			"err", errors.New("FileType错误"),
		)
		return errors.New("FileType错误")
	}

	if user.Type != 0 && user.Type != 1 {
		config.GetLogger().Warnw("Type错误",
			"err", errors.New("类型错误"),
		)
		return errors.New("类型错误")
	}

	if user.Type == 0 && user.MoveType == "abnormal" {
		config.GetLogger().Warnw("Type错误",
			"err", errors.New("类型错误"),
		)
		return errors.New("类型错误")
	}
	config.GetLogger().Info("进行数据有效性验证完成")

	//设置路径（绝对、相对）
	//dir := "D:\20study\2020project\back\"
	//.\data\tsyhh\origin\run\amp
	dir := ".\\data\\" + user.UserName + "\\" + key2[user.Type] + "\\" + key1[user.FileType-1] + "\\" + user.MoveType + "\\"
	filename := ""
	if user.Type == 0 {
		if user.MoveType == "amp" {
			filename = user.FileName + "_amp.json"
		} else if user.MoveType == "phase" {
			filename = user.FileName + "_phase.json"
		}
	} else if user.Type == 1 {
		if user.MoveType == "amp" {
			filename = "dealt_" + user.FileName + "_amp.json"
		} else if user.MoveType == "phase" {
			filename = "dealt_" + user.FileName + "_phase.json"
		} else if user.MoveType == "abnormal" {
			filename = "abnormal_" + user.FileName + ".json"
		}
	}

	fileStr := path.Join(dir, filename)
	//fileStr := dir + user.FileName

	fmt.Println(fileStr)

	//验证文件是否为json文件
	ext := path.Ext(fileStr)
	if ext != ".json" {
		config.GetLogger().Warnw("文件错误",
			"err", errors.New("文件扩展名不为.json"),
		)
		return errors.New("文件扩展名不为.json")
	}

	//读取文件
	f, errs := ioutil.ReadFile(fileStr)
	if errs != nil {
		fmt.Println("read fail", errs)
		config.GetLogger().Warnw("文件读取失败",
			"err", err,
		)
		return err
	}

	//读取json文件并传给前端
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
 * 根据用户名创建目录结构
 */
func createDir(userName string) (err error) {
	key1 := [3]string{"origin", "dealt", "upload"}
	key2 := [3]string{"run", "walk", "shake"}
	key3 := [3]string{"amp", "phase", "abnormal"}
	filepath := ""

	i := 0
	j := 0
	k := 0

	for i = 0; i < 3; i++ {
		for j = 0; j < 3; j++ {
			for k = 0; k < 3; k++ {
				if i == 0 {
					if k == 2 {
						continue
					}
					filepath = ".\\data\\" + userName + "\\" + key1[i] + "\\" + key2[j] + "\\" + key3[k]
				} else if i == 1 {
					filepath = ".\\data\\" + userName + "\\" + key1[i] + "\\" + key2[j] + "\\" + key3[k]
				} else {
					filepath = ".\\data\\" + userName + "\\" + key1[i] + "\\" + key2[j]
				}

				err = os.MkdirAll(filepath, os.ModePerm)
				if err != nil {
					config.GetLogger().Warnw("文件夹创建失败",
						"err", err,
					)
					return
				}
			}
		}
	}

	return
}

/**
 * 上传文件并解析
 */
func (u *User) upload(c *gin.Context) (err error) {
	data := &u.UploadData
	dataType := [3]string{"run", "walk", "shake"}
	dbName := [6]string{"origin_run", "origin_walk", "origin_shakehand", "dealt_run", "dealt_walk", "dealt_shakehand"}
	userID := ""

	config.GetLogger().Info("开始获取文件")
	//获取用户名
	name := c.PostForm("user_name")
	//获取文件类型：1是跑步，2是行走，3是摇手
	filetype, _ := strconv.Atoi(c.PostForm("type"))
	//获取文件
	file, err := c.FormFile("file")
	if err != nil {
		data.Uploaded = false
		config.GetLogger().Warnw("文件读取失败",
			"err", err,
		)
		return
	}
	if filetype < 1 || filetype > 3 {
		data.Uploaded = false
		config.GetLogger().Warnw("文件接收失败",
			"err", errors.New("文件类型错误"),
		)
		return errors.New("文件类型错误")
	}
	config.GetLogger().Info("获取文件结束")

	config.GetLogger().Info("开始文件重命名")
	//文件重命名
	ext := path.Ext(file.Filename)
	if ext != ".dat" {
		config.GetLogger().Warnw("文件错误",
			"err", errors.New("文件扩展名不为.dat"),
		)
		return errors.New("文件扩展名不为.dat")
	}

	file.Filename = dataType[filetype-1] + time.Now().Format("2006-01-02-15-04-05") + ext
	filename := dataType[filetype-1] + time.Now().Format("2006-01-02-15-04-05")
	fmt.Println(file.Filename)
	config.GetLogger().Info("文件重命名结束")

	config.GetLogger().Info("开始保存文件至服务器")
	filepath := ".\\data\\" + name + "\\upload\\" + dataType[filetype-1] + "\\"

	//设置保存路径
	dst := filepath + file.Filename

	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		data.Uploaded = false
		config.GetLogger().Warnw("文件上传失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("保存文件至服务器结束")

	config.GetLogger().Info("开始解析并生成原始数据")
	//设置运行时参数
	args := []string{"read_bfee_file.py", dataType[filetype-1] + time.Now().Format("2006-01-02-15-04-05"), name}
	//使用命令行的方式运行python文件
	cmd := exec.Command("python", args...)
	//设置执行文件的文件路径
	cmd.Dir = ".\\py\\"
	//开始执行
	err = cmd.Run()
	if err != nil {
		data.Uploaded = false
		config.GetLogger().Warnw("执行python脚本失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("解析并生成原始数据结束")

	config.GetLogger().Info("开始更新数据库")
	row := db.Table("user_info").Where("user = ?", name).Select("id").Row()
	err = row.Scan(&userID)
	if err != nil {
		data.Uploaded = false
		config.GetLogger().Warnw("查询用户ID失败",
			"err", err,
		)
		return
	}

	tx := db.Begin()
	od := new(OriginData)
	od.UID = userID
	od.FileName = file.Filename
	od.DataUrl = dst
	od.Time = time.Now().String()
	fmt.Println(dst)
	err = db.Table(dbName[filetype-1]).Create(&od).Error
	if err != nil {
		data.Uploaded = false
		config.GetLogger().Warnw("更新数据库失败",
			"err", err,
		)
		tx.Rollback()
		return
	}

	dd := new(DealtData)
	path1 := ".\\data\\" + name + "\\dealt\\" + dataType[filetype-1] + "\\amp\\" + "dealt_" + filename + "_amp.json"
	path2 := ".\\data\\" + name + "\\dealt\\" + dataType[filetype-1] + "\\phase\\" + "dealt_" + filename + "_phase.json"
	path3 := ".\\data\\" + name + "\\dealt\\" + dataType[filetype-1] + "\\abnormal\\" + "abnormal_" + filename + ".json"
	path4 := ".\\data\\" + name + "\\origin\\" + dataType[filetype-1] + "\\amp\\" + filename + "_amp.json"
	path5 := ".\\data\\" + name + "\\origin\\" + dataType[filetype-1] + "\\phase\\" + filename + "_phase.json"
	dd.UID = userID
	dd.FileName = filename
	dd.Amp = path1
	dd.Phase = path2
	dd.Abnormal = path3
	dd.OriginAmp = path4
	dd.OriginPhase = path5
	dd.Time = time.Now().String()
	err = db.Table(dbName[filetype+2]).Create(&dd).Error
	if err != nil {
		data.Uploaded = false
		config.GetLogger().Warnw("更新数据库失败",
			"err", err,
		)
		tx.Rollback()
		return
	}
	tx.Commit()
	config.GetLogger().Info("更新数据库结束")

	data.Uploaded = true

	return
}

/**
 * 统计信息
 */
func (u *User) statistics(cont string) (err error) {
	data := &u.StatisticsData

	//接收前端所传数据并解析
	config.GetLogger().Info("开始解析数据")
	id := new(ReceiveID)
	err = json.Unmarshal([]byte(cont), &id)
	if err != nil {
		config.GetLogger().Warnw("数据解析失败",
			"err", err,
		)
		return
	}
	config.GetLogger().Info("解析数据结束")

	config.GetLogger().Info("开始获取跑步数据条数")
	i := new(StatisticsData)
	err = db.Table("dealt_run").Where("id = ?", id).Count(&i.Value).Error
	if err != nil {
		config.GetLogger().Warnw("数据库错误",
			"err:", err,
		)
		return
	}
	i.Name = "Run"
	*data = append(*data, *i)
	config.GetLogger().Info("获取跑步数据条数结束")

	config.GetLogger().Info("开始获取行走数据条数")
	err = db.Table("dealt_walk").Where("id = ?", id).Count(&i.Value).Error
	if err != nil {
		config.GetLogger().Warnw("数据库错误",
			"err:", err,
		)
		return
	}
	i.Name = "Walk"
	*data = append(*data, *i)
	config.GetLogger().Info("获取行走数据条数结束")

	config.GetLogger().Info("开始获取摇手数据条数")
	err = db.Table("dealt_shakehand").Where("id = ?", id).Count(&i.Value).Error
	if err != nil {
		config.GetLogger().Warnw("数据库错误",
			"err:", err,
		)
		return
	}
	i.Name = "ShakeHands"
	*data = append(*data, *i)
	config.GetLogger().Info("获取摇手数据条数结束")

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

func (u *User) GetMovementListData(cont string) (err error, data MovementListData) {
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

func (u *User) GetStatisticsData(cont string) (err error, data []StatisticsData) {
	config.GetLogger().Info("开始获取统计数据")

	err = u.statistics(cont)

	data = u.StatisticsData

	config.GetLogger().Info("获取统计数据结束")

	return
}

func (u *User) GetUploadData(c *gin.Context) (err error, data UploadData) {
	config.GetLogger().Info("开始上传文件")

	err = u.upload(c)

	data = u.UploadData

	config.GetLogger().Info("上传文件结束")

	return
}
