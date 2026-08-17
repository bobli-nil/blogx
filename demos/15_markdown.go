package main

import (
	"blogx_server/utils/markdown"
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var str = `
# 这是一级标题

> 这是引用

![abc](test.jpg)
`

func main() {
	htmlStr := markdown.MdToHTML(str)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlStr)))
	if err != nil {
		fmt.Println(err)
		return
	}
	htmlText := doc.Text()
	fmt.Println(htmlText)
}
