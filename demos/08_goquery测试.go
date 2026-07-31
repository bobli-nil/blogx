package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	file, _ := os.Open("demos/index.html")
	defer file.Close()
	//fmt.Println(file.Name(), err)
	doc, _ := goquery.NewDocumentFromReader(file)

	s := doc.Find("link[rel='icon']")
	fmt.Println(s.Length())

	//head := doc.Find("head")
	//head.AppendHtml("<title>这是标题</title>")
	//
	//html, _ := doc.Html()
	//
	//os.WriteFile("demos/index.html", []byte(html), 0777)
}
