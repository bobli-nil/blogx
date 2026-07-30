package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"fmt"
)

func main() {
	flags.Parse()
	global.Conf = core.ReadConf()

	token, err := jwt.GenerateToken(1, "张三", enum.AdminRole)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

	//t := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsInVzZXJOYW1lIjoi5byg5LiJIiwicm9sZSI6MSwiaXNzIjoibG1yIiwiZXhwIjoxNzg1Mzg3NjE1LCJuYmYiOjE3ODUzODc2MTQsImlhdCI6MTc4NTM4NzYxNH0.ajv4Ea0gDS10yZVjCrivySkL-YWGz4lftSC13wPYQ4I"
	//
	//claims, err := jwt.ParseToken(t)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(claims)
}
