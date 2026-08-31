package main

import (
	"blogx_server/models"
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

	list := text_service.MdContentTransformation(models.ArticleModel{
		Model: models.Model{
			ID: 1,
		},
		Title:   "gvb博客功能开发",
		Content: string(byteData),
	})
	fmt.Printf("%q \n", list)
}
