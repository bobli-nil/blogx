package main

import (
	"fmt"
	"reflect"
)

type Model struct {
	ID uint
}

type User struct {
	Model
	name string
}

func main() {
	a := User{
		Model: Model{
			ID: 1,
		},
		name: "lisi",
	}
	fmt.Println(a.ID)
	of := reflect.TypeOf(a)
	name, b := of.FieldByName("ID")
	fmt.Println(name, b)
}
