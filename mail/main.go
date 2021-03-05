package mail

import (
	"SHU/parseSetting"

	"gopkg.in/gomail.v2"
)

/*
go邮件发送
*/

func SendMail(sender,authorization string,mailTo []string, subject string, body string) error {
	// 设置邮箱主体
	mailConn := map[string]string{
		"user": sender, //发送人邮箱（邮箱以自己的为准）
		"pass": authorization, //发送人邮箱授权码 开启授权密码后在pass填写授权码
		"host": "smtp.qq.com",      //邮箱服务器（此时用的是qq邮箱）
	}

	m := gomail.NewMessage(
		//发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		//就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "LLL")) // 添加别名
	m.SetHeader("To", mailTo...)                                  // 发送给用户(可以多个)
	m.SetHeader("Subject", subject)                               // 设置邮件主题
	m.SetBody("text/html", body)                                  // 设置邮件正文

	/*
	   创建SMTP客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为465的话，
	   自动开启SSL，这个时候需要指定TLSConfig
	*/
	d := gomail.NewDialer(mailConn["host"], 465, mailConn["user"], mailConn["pass"]) // 设置邮件正文
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}

func SendErro(info string, trac string, setting parseSetting.Setting) error {
	// 邮件接收方
	mailTo := setting.Mail.Receivers

	subject := info // 邮件主题
	body := trac    // 邮件正文

	err := SendMail(setting.Mail.Sender,setting.Mail.AuthorizationCode,mailTo, subject, body)
	return err
}
