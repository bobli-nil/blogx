package main

import (
	"blogx_server/service/text_service"
	"fmt"
	"os"
)

func main() {
	byteData, err := os.ReadFile("demos/md/text.md")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	list := text_service.MdContentTransformation(1, "", string(byteData))
	fmt.Printf("%q \n", list)
}
