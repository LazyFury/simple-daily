package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Mail 邮件配置
type Mail struct {
	Nickname string `json:"nickname"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// Auth 身份认证
func (m *Mail) Auth() smtp.Auth {
	return smtp.PlainAuth("", m.User, m.Password, m.Host)
}

// SendMail 发送邮件
func (m *Mail) SendMail(subject string, to []string, body string) (err error) {

	template := `To:%s
From:%s<%s>
Subject:%s
Content-Type: text/html; charset=UTF-8

<html>
	<body>
		%s
	<body>
</html>
	`

	msg := []byte(fmt.Sprintf(template, strings.Join(to, ","), m.Nickname, m.User, subject, body))

	return smtp.SendMail(m.Host+":"+m.Port, m.Auth(), m.User, to, msg)
}
