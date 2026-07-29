package main

import (
	"blogx_server/utils/jwt"
	"fmt"
)

func main() {
	//token, err := jwt.GenerateToken(18, "lisi")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(token)

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjE4LCJ1c2VyTmFtZSI6Imxpc2kiLCJpc3MiOiJsbXIiLCJleHAiOjE3ODUzNzg2ODEsIm5iZiI6MTc4NTI5MjI4MSwiaWF0IjoxNzg1MjkyMjgxfQ.pcTG6TwANU_w8gEqbGJZZ5R4XRa4yka7UdmtVRHPl8Y"
	claims, err := jwt.ParseToken(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", claims)
	fmt.Println(claims.UserID, claims.UserName)
	fmt.Printf("%T\n", claims.RegisteredClaims.ExpiresAt)
}
