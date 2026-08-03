package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/email_service"
	"fmt"
)

func main() {
	//m := gomail.NewMessage()
	//m.SetHeader("From", "1172778182@qq.com")
	//m.SetHeader("To", "18813070947@163.com")
	//m.SetHeader("Subject", "邮件测试")
	//m.SetBody("text/html", "Hello world")
	//
	//d := gomail.NewDialer("smtp.qq.com", 587, "1172778182@qq.com", "cvjtoycgcjyvfegh")
	//
	//if err := d.DialAndSend(m); err != nil {
	//	panic(err)
	//}

	flags.Parse()
	global.Conf = core.ReadConf()

	err := email_service.SendRegisterCode("18813070947@163.com", "1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("邮件发送成功")
}
