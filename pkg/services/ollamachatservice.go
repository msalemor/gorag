package services

import (
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

type OllamaChatServiceOptions struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type OllamaChatService struct {
	Client   *http.Client
	Endpoint string
	Model    string
}

type OllamaChatRequest struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	Messages []Message
	Options  OllamaChatServiceOptions `json:"options"`
	Stream   bool                     `json:"stream"`
}

type OllamaChatResponse struct {
	Model              string  `json:"model"`
	CreateAt           string  `json:"created_at"`
	Message            Message `json:"message"`
	Done               bool    `json:"done"`
	TotalDuration      int     `json:"total_duration"`
	LoadDuration       int     `json:"load_duration"`
	PromptEvalCount    int     `json:"prompt_eval_count"`
	PromptEvalDuration int     `json:"prompt_eval_duration"`
	EvalCount          int     `json:"eval_count"`
	EvalDuration       int     `json:"eval_duration"`
}

func (e *OllamaChatService) Chat(opts *ChatOpts) *OllamaChatResponse {

	// marshall data to json (like json_encode)
	ollamaReq := OllamaChatRequest{
		Model:    e.Model,
		Messages: opts.Messages,
		Options: OllamaChatServiceOptions{
			Temperature: opts.Temperature,
			MaxTokens:   opts.MaxTokens,
		},
		Stream: opts.Stream,
	}

	marshalled, err := json.Marshal(ollamaReq)
	if err != nil {
		logrus.Error(fmt.Sprintf("impossible to marshall teacher: %s", err))
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, e.Endpoint, bytes.NewReader(marshalled))
	if err != nil {
		logrus.Error(fmt.Sprintf("impossible to create request: %s", err))
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", "ollama")

	resp, err := e.Client.Do(req)
	if err != nil {
		logrus.Error(fmt.Sprintf("impossible to send request: %s", err))
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(fmt.Sprintf("impossible to read response body: %s", err))
		return nil
	}

	// Parse the response body into a vector of floats
	var completion OllamaChatResponse
	err = json.Unmarshal(body, &completion)
	if err != nil {
		logrus.Error(fmt.Sprintf("impossible to unmarshal response body: %s", err))
		return nil
	}

	return &completion
}
