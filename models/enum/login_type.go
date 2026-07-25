package enum

type LoginType uint8

const (
	UserNamePwdLoginType LoginType = 1 // 用户名密码登录
	QQLoginType          LoginType = 2 // QQ登录
	EmailLoginType       LoginType = 3 // 邮箱登录
)
