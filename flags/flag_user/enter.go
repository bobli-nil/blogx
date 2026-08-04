package flag_user

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/pwd"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type User struct{}

type Option struct {
	value enum.RoleType
	label string
}

func (User) Create() {
	// 输入角色、校验角色
	var roleName string
	var role enum.RoleType
	survey.AskOne(&survey.Select{
		Message: "请选择角色",
		Options: []string{"超级管理员", "普通用户", "访客"},
	}, &roleName)

	switch roleName {
	case "超级管理员":
		role = enum.AdminRole
		break
	case "普通用户":
		role = enum.UserRole
		break
	case "访客":
		role = enum.VisitorRole
		break
	}

	// 输入用户名、检查是否存在
	var username string
	survey.AskOne(&survey.Input{
		Message: "请输入用户名",
	}, &username)
	err := global.DB.Take(&models.UserModel{}, "username = ?", username).Error
	if err == nil {
		fmt.Println("该用户名已存在")
		return
	}

	// 输入密码、特殊手段，输入两次密码，检查两次是否一样，加密
	var password string
	survey.AskOne(&survey.Password{
		Message: "请输入密码:",
	}, &password)

	var rePassword string
	survey.AskOne(&survey.Password{
		Message: "请再次输入密码:",
	}, &rePassword)

	if password != rePassword {
		fmt.Println("密码不一致")
		return
	}

	// 创建用户
	passwordHash, _ := pwd.GenerateFromPassword(password)
	err = global.DB.Create(&models.UserModel{
		Username: username,
		Nickname: username,
		Password: passwordHash,
		Role:     role,
	}).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("创建成功")
}
