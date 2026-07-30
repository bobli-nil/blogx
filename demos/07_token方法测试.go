package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/utils/jwt"
	"fmt"
)

func main() {
	flags.Parse()
	global.Conf = core.ReadConf()

	//token, err := jwt.GenerateToken(1, "张三", enum.AdminRole)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(token)

	t := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoi5byg5LiJIiwicm9sZSI6MSwiaXNzIjoibG1yIiwiZXhwIjoxNzg1NDkwOTA2LCJuYmYiOjE3ODU0MDQ1MDYsImlhdCI6MTc4NTQwNDUwNn0.0BjP-NoF4i_WjnqichI6U_bORRpxRaCI2xfl4NWLAQM"
	claims, err := jwt.ParseToken(t)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", claims)
}
