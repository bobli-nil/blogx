package text_service

import (
	"fmt"
	"strings"
)

type TextModel struct {
	ArticleID uint   `json:"article_id"`
	Head      string `json:"head"`
	Body      string `json:"body"`
}

func MdContentTransformation(id uint, title, content string) (list []TextModel) {
	lines := strings.Split(content, "\n")

	var headList []string
	var bodyList []string
	var body string
	headList = append(headList, title)
	var flag bool
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			flag = !flag
		}
		// 处理标题行和上一个标题的body
		if !flag && strings.HasPrefix(line, "#") {
			bodyList = append(bodyList, GetBody(body))
			body = ""
			headList = append(headList, GetHead(line))
			continue
		}
		body += line
	}
	if body != "" {
		bodyList = append(bodyList, GetBody(body))
	}

	// 单独处理以标题结尾的情况
	if len(headList) > len(bodyList) {
		bodyList = append(bodyList, "")
	}

	if len(headList) != len(bodyList) {
		fmt.Println("headList和bodyList数量不一致")
		fmt.Printf("%q %d \n", headList, len(headList))
		fmt.Printf("%q %d \n", bodyList, len(bodyList))
		return
	}

	for i := 0; i < len(headList); i++ {
		list = append(list, TextModel{
			ArticleID: id,
			Head:      headList[i],
			Body:      bodyList[i],
		})
	}

	return

}

func GetHead(head string) string {
	return strings.TrimSpace(strings.Join(strings.Split(head, " ")[1:], " "))
}

func GetBody(body string) string {
	return strings.TrimSpace(body)
}
