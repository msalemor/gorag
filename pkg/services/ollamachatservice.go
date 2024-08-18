package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func (e *OllamaChatService) Chat(messages []Message, temperature float64, maxTokens int, stream bool) *OllamaChatResponse {

	// marshall data to json (like json_encode)
	ollamaReq := OllamaChatRequest{
		Model:    e.Model,
		Messages: messages,
		Options: OllamaChatServiceOptions{
			Temperature: temperature,
			MaxTokens:   maxTokens,
		},
		Stream: stream,
	}

	marshalled, err := json.Marshal(ollamaReq)
	if err != nil {
		log.Printf("impossible to marshall teacher: %s", err)
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, e.Endpoint, bytes.NewReader(marshalled))
	if err != nil {
		log.Printf("impossible to create request: %s", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", "ollama")

	resp, err := e.Client.Do(req)
	if err != nil {
		log.Printf("impossible to send request: %s", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("impossible to read response body: %s", err)
		return nil
	}

	// Parse the response body into a vector of floats
	var completion OllamaChatResponse
	err = json.Unmarshal(body, &completion)
	if err != nil {
		log.Printf("impossible to unmarshal response body: %s", err)
		return nil
	}

	return &completion
}
