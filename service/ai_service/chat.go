package ai_service

import (
	"encoding/json"
	"io"

	"github.com/sirupsen/logrus"
)

func Chat(content string) (msg string, err error) {
	r := Request{
		Model: "gpt-5.4-mini",
		Messages: []Message{
			Message{Role: "system", Content: "AI"},
			Message{Role: "user", Content: content},
		},
	}
	res, err := BaseRequest(r)
	defer res.Body.Close()
	if err != nil {
		logrus.Error(err)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return
	}

	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Errorf("解析失败 %s %s", err, string(body))
		return
	}
	msg = response.Choices[0].Message.Content
	return
}
