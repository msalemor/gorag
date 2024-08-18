package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show configured endpoints and models",
		Long:  "Show configured endpoints and models",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Configured models and endpoints:\n\n")
			fmt.Println("OllamaChatEndpoint:", OllamaChatEndpoint)
			fmt.Println("ChatModel:", ChatModel)
			fmt.Println("")
			fmt.Println("OllamaEmbeddingEndpoint:", OllamaEmbeddingEndpoint)
			fmt.Println("EmbeddingModel:", EmbeddingModel)
		},
	}
}
