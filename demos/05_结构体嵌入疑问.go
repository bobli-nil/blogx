package main

import (
	"encoding/json"
	"fmt"
)

type Log struct {
	User User
	ID   uint
}

type User struct {
	UserID uint
	Name   string
}

func main() {
	log := &Log{
		User: User{
			UserID: 1,
			Name:   "lisi",
		},
		ID: 1,
	}
	marshal, _ := json.Marshal(log)
	fmt.Println(string(marshal))
}
