# SHU_REPOT_GO

![avatar](https://img.shields.io/badge/Async-Yes-red)
![avatar](https://img.shields.io/badge/Configuration-Yes-green)
![avatar](https://img.shields.io/badge/license-MIT-blue)

Go语言实现的上海大学每日一报项目, 摆脱selenium纯api封装, 
错误自动发送邮件通知
## 项目结构
```
SHU_report_public
│
├── go.mod
├── go.sum
│
├── mail
│   └── main.go
│
├── parseSetting
│   └── main.go
│
├── shu_f
│   └── server.py
│
├── report
│   ├── report.go
│   └── report_test.go
│
├── main.go
├── readme.md
└── setting.json
```

## 更新改动
(21-4-22更):
pyppeteer 现在只负责js加密部分的运行, 依旧保持CS架构, 但是现在无头浏览器只负责运行js, 其他请求全部代码实现, 性能和稳定性较快.

(21-4-21更):
使用了CS架构 python(fastapi+pyppeteer) + golang的结构

但是这样的话还不如直接完全使用python实现来得方便

## 启动方法(目前启动部署比较麻烦了, 就不更了)
先启动server
```python
# 在shu_f文件夹下的server, 用于获取cookie
if __name__ == '__main__':
    uvicorn.run(app, host="0.0.0.0", port=8989)

# 启动完成之后, 需要在终端下使用curl 127.0.0.1:8989/init来完成browser的装载进入内存
# 服务器基本都是靠终端交互, 要避免浏览器被杀进程,这部分比较麻烦, 后续更新解决方案, 感兴趣的朋友可以自己动手尝试
```
然后启动client, 负责定时报送和发送错误邮件, 注意自己的路径
```shell
go run "./main.go" 
```

## ios快捷指令来查看服务和报送状态
ios可以通过快捷指令来快捷查看报送状态


## 使用方法(旧版本)

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

