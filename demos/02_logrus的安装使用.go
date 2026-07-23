package main

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

var Logger = new(logrus.Entry)

func main() {
	//setColor()
	useBytesBuffer()
}

func useLogrus() {
	logrus.SetLevel(logrus.DebugLevel)
	fmt.Println(logrus.GetLevel())
	fmt.Println("-----------")

	logrus.Debugln("Debugln")
	logrus.Infoln("Infoln")
	logrus.Warnln("Warnln")
	logrus.Errorln("Errorln")
	logrus.Println("Println")

	logrus1 := logrus.WithField("project", "study")
	logrus1.Infoln("logrus1")
	logrus2 := logrus.WithFields(logrus.Fields{"a1": "v1"})
	logrus2.Infoln("logrus2")
}

func initLogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	Logger = logrus.WithFields(logrus.Fields{
		"request_id": "id",
		"name":       "lisi",
	})
	Logger.Infoln("hello world")
	Logger.Debugln("hello world")
}

func setColor() {
	// 前景色
	fmt.Println("\033[30m 黑色 \033[0m")
	fmt.Println("\033[31m 红色 \033[0m")
	fmt.Println("\033[32m 绿色 \033[0m")
	fmt.Println("\033[33m 黄色 \033[0m")
	fmt.Println("\033[34m 蓝色 \033[0m")
	fmt.Println("\033[35m 紫色 \033[0m")
	fmt.Println("\033[36m 青色 \033[0m")
	fmt.Println("\033[37m 灰色 \033[0m")
	// 背景色
	fmt.Println("\033[40m 黑色 \033[0m")
	fmt.Println("\033[41m 红色 \033[0m")
	fmt.Println("\033[42m 绿色 \033[0m")
	fmt.Println("\033[43m 黄色 \033[0m")
	fmt.Println("\033[44m 蓝色 \033[0m")
	fmt.Println("\033[45m 紫色 \033[0m")
	fmt.Println("\033[46m 青色 \033[0m")
	fmt.Println("\033[47m 灰色 \033[0m")
}

func useBytesBuffer() {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "hello world, ")
	fmt.Fprintf(&buf, "我的名字是%s \n", "里斯")
	fmt.Fprintf(&buf, "这是另一行")

	fmt.Println(string(buf.Bytes()))
}
