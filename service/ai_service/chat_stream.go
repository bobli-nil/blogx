package ai_service

import (
	"bufio"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Choice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Logprobs     any `json:"logprobs"`
	FinishReason any `json:"finishReason"`
}
type StreamData struct {
	Id                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Object            string   `json:"object"`
	SystemFingerprint any      `json:"system_fingerprint"`
}

func ChatStream(content string) (msgChan chan string, err error) {
	msgChan = make(chan string)
	r := Request{
		Model: "gpt-5.4-mini",
		Messages: []Message{
			Message{Role: "system", Content: "AI"},
			Message{Role: "user", Content: content},
		},
		Stream: true,
	}
	res, err := BaseRequest(r)
	if err != nil {
		logrus.Error(err)
		return
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Split(bufio.ScanLines)

	go func() {
		defer res.Body.Close()
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}
			data := text[6:]
			if data == "[DONE]" {
				close(msgChan)
				return
			}
			var item StreamData
			err = json.Unmarshal([]byte(data), &item)
			if err != nil {
				logrus.Errorf("解析失败 %s %s", err, text)
				continue
			}
			if len(item.Choices) > 0 {
				msgChan <- item.Choices[0].Delta.Content
			}
		}
	}()

	return
}
