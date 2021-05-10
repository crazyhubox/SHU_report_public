# SHU_REPOT_GO

![avatar](https://img.shields.io/badge/Async-Yes-red)
![avatar](https://img.shields.io/badge/Configuration-Yes-green)
![avatar](https://img.shields.io/badge/license-MIT-blue)

Go语言实现的上海大学每日一报项目, 摆脱selenium纯api封装, 错误自动发送邮件通知

## 项目结构

```bash
SHU_report_public
│
├── LICENSE
├── README.md
├── go.mod
├── go.sum
├── js_test
│   └── jiami.js
├── mail
│   └── main.go
├── parseSetting
│   └── main.go
├── report
│   ├── cookie.go
│   ├── pdatas.go
│   ├── report.go
│   ├── report_test.go
│   └── result.go
├── shu_f
│   ├── __init__.py
│   ├── aync_test.py
│   ├── cookie.py
│   ├── pypp_cookie.py
│   ├── server.py
│   ├── serverlog.log
│   ├── start.sh
│   ├── stop.sh
│   └── test.py
├── setting.json
├── main.go
├── makefile
└── static
    └── check.gif
```

## 更新改动

(21-5-10更):
pyppeteer启动的无头浏览器仅作为一个javascript的运行容器,并且作为一个服务在本地后台等待.将获取加密参数的过程封装成一个rpc过程

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

![image](https://github.com/crazyhubox/SHU_report_public/blob/main/static/check.gif)

## 使用方法(旧版本)

1. 安装golang环境

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

## 最后

项目仅仅作为同学们参考和讨论交流, 并不属于拿来即用类型的项目,如果需要实际使用还需要根据自身情况进行修改.