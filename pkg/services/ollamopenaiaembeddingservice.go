package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
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

type OllamaOpenAIEmbeddingData struct {
	Ojbect    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
}

type OllamaOpenAIEmbedResponse struct {
	Object string                      `json:"object"`
	Data   []OllamaOpenAIEmbeddingData `json:"data"`
}

func (e *OllamaOpenAIEmbeddingService) Embed(opts *EmbeddingOpts) *[]float64 {

	// marshall data to json (like json_encode)
	ollamaReq := OllamaOpenAIEmbedRequest{
		Model: e.Model,
		Input: opts.Text,
	}

	marshalled, err := json.Marshal(ollamaReq)
	if err != nil {
		logrus.Error(fmt.Sprintf("Unable to marshall OllamaEmbedRequest: %s", err))
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, e.Endpoint, bytes.NewReader(marshalled))
	if err != nil {
		logrus.Error(fmt.Sprintf("Unable to create a NewRequest: %s", err))
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", "ollama")

	resp, err := e.Client.Do(req)
	if err != nil {
		logrus.Error(fmt.Sprintf("Unable to send the request: %s", err))
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(fmt.Sprintf("Unable to read response body: %s", err))
		return nil
	}

	// Parse the response body into a vector of floats
	var vector OllamaOpenAIEmbedResponse
	err = json.Unmarshal(body, &vector)
	if err != nil {
		logrus.Error(fmt.Sprintf("Unable to unmarshal response body to OllamaOpenAIEmbedResponse: %s", err))
		return nil
	}

	return &vector.Data[0].Embedding
}
