package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type OllamaOpenAIChatService struct {
	Client   *http.Client
	Endpoint string
	Model    string
}

type OllamaOpenAIChatRequest struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Messages    []Message
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

type OllamaOpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (e *OllamaOpenAIChatService) Chat(messages []Message, temperature float64, maxTokens int, stream bool) *OllamaOpenAIChatResponse {

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
	var completion OllamaOpenAIChatResponse
	err = json.Unmarshal(body, &completion)
	if err != nil {
		log.Printf("impossible to unmarshal response body: %s", err)
		return nil
	}

	return &completion
}
