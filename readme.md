# WiFi信号识别人体动作项目
---
## 技术
本次项目我们选择Go语言进行后端的开发，选用Ginroute框架、Logrus日志框架、Gorm数据库映射，以后也可以以此为模版进行新项目的开发。
## 目录说明
``` 
.
|-app
| |-config # 配置文件读取文件
|   |-logger.go # 日志配置
|   |-main.go
未完待续...
``` 
## 运行方法 配置文件配置项说明
运行时需要加入命令行参数，读取配置文件相关命令如下：
``` 
./main ./dev.yml
``` 
配置文件说明如下：
``` 
database:
  address: 118.31.171.61  #数据库域名
  port: 3306              #数据库端口
  dbname: project         #数据库名称
  user: project           #用户名
  password: WiFi6666      #数据库密码

```
## API文档
### 用户登录
* URL: /login
* Method: POST
#### 前端发送
``` 
{
    user_name: string
    user_pwd:  string
}

``` 
#### 返回数据
``` 
{
    status: "success" || "error"
    message: string
    data: {
        is_login: boolean //true表示登录成功
        type：    int//权限0是普通用户，1是管理员，2是超级管理员。
    }
}

``` 
