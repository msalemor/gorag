package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type OllamaEmbeddingService struct {
	Endpoint   string
	Model      string
	Dimensions int
	Client     *http.Client
}

type OllamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaEmbedResponse struct {
	Embedding []float64 `json:"embedding"`
}

func (e *OllamaEmbeddingService) Embed(text string) *[]float64 {

	// marshall data to json (like json_encode)
	ollamaReq := OllamaEmbedRequest{
		Model:  e.Model,
		Prompt: text,
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
	var vector OllamaEmbedResponse
	err = json.Unmarshal(body, &vector)
	if err != nil {
		log.Printf("Unable to unmarshal response body to OllamaEmbedResponse: %s", err)
		return nil
	}

	return &vector.Embedding
}