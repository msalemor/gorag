package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_COLLECTION = "FAQ"
)

var (
	Collection              = "FAQ"
	OllamaChatEndpoint      = "http://localhost:11434/v1/chat/completions"
	ChatModel               = "llama3"
	OllamaEmbeddingEndpoint = "http://localhost:11434/v1/embeddings"
	EmbeddingModel          = "nomic-embed-text"
	Verbose                 = false
	Keep                    = false
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		println("Error getting executable path")
	}
	exPath := filepath.Dir(ex)
	log.Println("Executable path: ", exPath)

	godotenv.Load(".env", exPath+"/.env")
	if value := os.Getenv("APP_COLLECTION"); value != "" {
		Collection = value
	}
	if value := os.Getenv("OLLAMA_CHAT_ENDPOINT"); value != "" {
		OllamaChatEndpoint = value
	}
	if value := os.Getenv("CHAT_MODEL"); value != "" {
		ChatModel = value
	}
	if value := os.Getenv("OLLAMA_EMBBEDING_ENDPOINT"); value != "" {
		OllamaEmbeddingEndpoint = value
	}
	if value := os.Getenv("EMBBEDING_MODEL"); value != "" {
		EmbeddingModel = value
	}
}
