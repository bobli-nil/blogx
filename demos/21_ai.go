package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/ai_service"
	"fmt"
	"io"
	"net/http"
)

func main() {
	flags.Parse()
	global.Conf = core.ReadConf()
	//GeModelList()
	msg, err := ai_service.Chat("你能帮我做什么")
	fmt.Println(msg, err)
}

func GeModelList() {
	url := "https://api.chatanywhere.tech/v1/models"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Conf.Ai.SecretKey))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
