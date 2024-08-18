package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const (
	DEFAULT_COLLECTION = "FAQ"
)

var (
	Collection              = "eChampShop"
	OllamaEmbeddingEndpoint = "http://localhost:11434/api/embeddings"
	ChatModel               = "llama3"
	OllamaChatEndpoint      = "http://localhost:11434/api/chat"
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

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gorag",
		Short: "gorag is a package and a CLI tool",
		Long:  "gorag is a package and a CLI tool long description",
	}

	root.AddCommand(consoleCmd())
	root.AddCommand(uiCmd())
	root.AddCommand(verCmd())
	root.AddCommand(configCmd())

	return root
}
