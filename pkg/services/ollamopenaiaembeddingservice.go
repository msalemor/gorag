package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type OllamaOpenAIEmbeddingService struct {
	Endpoint   string
	Model      string
	Dimensions int
	Client     *http.Client
}

type OllamaOpenAIEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type OllamaOpenAIEmbedResponse struct {
	Object string `json:"object"`
	Data   struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func (e *OllamaOpenAIEmbeddingService) Embed(text string) *[]float64 {

	// marshall data to json (like json_encode)
	ollamaReq := OllamaOpenAIEmbedRequest{
		Model: e.Model,
		Input: text,
	}

	marshalled, err := json.Marshal(ollamaReq)
	if err != nil {
		log.Printf("Unable to marshall OllamaEmbedRequest: %s", err)
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, e.Endpoint, bytes.NewReader(marshalled))
	if err != nil {
		log.Printf("Unable to create a NewRequest: %s", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", "ollama")

	resp, err := e.Client.Do(req)
	if err != nil {
		log.Printf("Unable to send the request: %s", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read response body: %s", err)
		return nil
	}

	// Parse the response body into a vector of floats
	var vector OllamaOpenAIEmbedResponse
	err = json.Unmarshal(body, &vector)
	if err != nil {
		log.Printf("Unable to unmarshal response body to OllamaEmbedResponse: %s", err)
		return nil
	}

	return &vector.Data.Embedding
}
