package emailsend

import (
	"crypto/tls"
	"gf_chat_server/utility/envtool"

	"gopkg.in/gomail.v2"
)

type EmailSendIns struct {
	emailIns *gomail.Dialer
}

func New() EmailSendIns {
	host := envtool.GetEnvValue("host", "smtp.qq.com").(string)
	port := envtool.GetEnvValue("port", 25).(int)
	userName := envtool.GetEnvValue("username", "").(string)
	password := envtool.GetEnvValue("password", "").(string) // qq邮箱填授权码
	ins := gomail.NewDialer(
		host,
		port,
		userName,
		password,
	)
	ins.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 关闭TLS认证

	return EmailSendIns{emailIns: ins}
}

// 封装消息
func (e *EmailSendIns) Message(addr string) *gomail.Message {
	m := gomail.NewMessage(
		gomail.SetEncoding(gomail.Base64),
	)
	m.SetHeader("From", e.emailIns.Username) // 发件人
	m.SetHeader("To", addr)
	return m
}

// 发送文字邮件
func (e *EmailSendIns) Send(addr string, str string) bool {
	msg := e.Message(addr)
	msg.SetBody("text/plain", str)
	if err := e.emailIns.DialAndSend(msg); err != nil {
		panic(err)
	}
	return false
}

// 发送HTML邮件
func (e *EmailSendIns) SendHTML(addr string, str string) bool {
	msg := e.Message(addr)
	msg.SetBody("text/html", str)
	if err := e.emailIns.DialAndSend(msg); err != nil {
		panic(err)
	}
	return false
}
