package ai_service

import (
	"blogx_server/global"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Id      string `json:"id"`
	Choices []struct {
		Message struct {
			Role        string        `json:"role"`
			Content     string        `json:"content"`
			Refusal     interface{}   `json:"refusal"`
			Annotations []interface{} `json:"annotations"`
		} `json:"message"`
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens            int `json:"prompt_tokens"`
		CompletionTokens        int `json:"completion_tokens"`
		TotalTokens             int `json:"total_tokens"`
		CompletionTokensDetails struct {
			AudioTokens     int `json:"audio_tokens"`
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
		PromptTokensDetails struct {
			AudioTokens  int `json:"audio_tokens"`
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
	} `json:"usage"`
	Routing struct {
		ServingPipereplica string `json:"serving_pipereplica"`
	} `json:"routing"`
	ServiceTier string `json:"service_tier"`
}

const baseUrl = "https://api.chatanywhere.tech/v1/chat/completions"

func Chat(content string) (msg string, err error) {
	method := "POST"
	r := Request{
		Model: "gpt-5.4-mini",
		Messages: []Message{
			Message{Role: "system", Content: "你是blogx的文章分析AI"},
			Message{Role: "user", Content: content},
		},
	}
	byteData, _ := json.Marshal(r)
	req, err := http.NewRequest(method, baseUrl, bytes.NewBuffer(byteData))
	if err != nil {
		logrus.Error(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Conf.Ai.SecretKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("请求失败 %s", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
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
