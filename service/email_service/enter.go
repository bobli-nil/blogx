package email_service

import (
	"blogx_server/global"
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendRegisterCode(to, code string) error {
	em := global.Conf.Email
	subject := fmt.Sprintf("【%s】账号注册", em.SendNickname)
	body := fmt.Sprintf("你正在进行账号注册，验证码 %s ，十分钟内有效", code)
	return SendEmail([]string{to}, subject, body)
}

func SendResetPwdCode(to, code string) error {
	em := global.Conf.Email
	subject := fmt.Sprintf("【%s】密码重置", em.SendNickname)
	body := fmt.Sprintf("你正在进行密码重置，验证码 %s ，十分钟内有效", code)
	return SendEmail([]string{to}, subject, body)
}

func SendBindEmailCode(to, code string) error {
	em := global.Conf.Email
	subject := fmt.Sprintf("【%s】邮箱绑定", em.SendNickname)
	body := fmt.Sprintf("你正在进行邮箱绑定，验证码 %s ，十分钟内有效", code)
	return SendEmail([]string{to}, subject, body)
}

func SendEmail(to []string, subject, body string) error {
	em := global.Conf.Email
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(em.SendEmail, em.SendNickname))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(em.Domain, em.Port, em.SendEmail, em.AuthCode)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
