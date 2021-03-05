# SHU_REPOT_GO
Go语言实现的上海大学每日一报项目, 摆脱selenium纯api封装, 
错误自动发送邮件通知
## 项目结构
```
SHU_report_public

├── go.mod
├── go.sum
├── mail
│   └── main.go
├── main.go
├── parseSetting
│   └── main.go
├── readme.md
├── report
│   ├── report.go
│   └── report_test.go
└── setting.json
```

## 使用方法

1. 安装golang环境

安装初始化项目

``` shell
go mod tidy
go mod download
go mod vendor
```

2. 配置setting.json

``` json
{
    "mail":{
        "sender":"qq邮箱@qq.com",
        "receivers":["qq邮箱@qq.com"],
        "authorization_code":"邮箱授权码"
    },
    "report_info":{
        "stu_num":"学号",
        "password":"密码"
    },
    "runtime_cron":"30 6 * * */1 这个表示每天的6点半执行,去掉中文"
}
```

3. 运行项目

``` 
go run main.go
```
