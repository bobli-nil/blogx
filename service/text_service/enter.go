package text_service

import (
	"blogx_server/models"
	"fmt"
	"strings"
)

func MdContentTransformation(model models.ArticleModel) (list []models.TextModel) {
	lines := strings.Split(model.Content, "\n")

	var headList []string
	var bodyList []string
	var body string
	headList = append(headList, model.Title)
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

	if len(headList) != len(bodyList) {
		fmt.Println("headList和bodyList数量不一致")
		fmt.Printf("%q %d \n", headList, len(headList))
		fmt.Printf("%q %d \n", bodyList, len(bodyList))
		return
	}

	for i := 0; i < len(headList); i++ {
		list = append(list, models.TextModel{
			ArticleID: model.ID,
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
