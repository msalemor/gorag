package cmd

import (
	"fmt"
	"log"

	"github.com/msalemor/gorag/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//https://www.fdic.gov/system/files/2024-07/banklist.csv

func processUI() {
	r := process.ConfigureRoutes(OllamaChatEndpoint, OllamaEmbeddingEndpoint, DEFAULT_COLLECTION, EmbeddingModel, false, Verbose)
	log.Println("Starting server on http://localhost:8080")
	r.Run()
}

func uiCmd() *cobra.Command {
	cmdUI := &cobra.Command{
		Use:   "ui",
		Short: "Run in UI mode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running in UI mode")
			fmt.Println("Verbose:", Verbose)
			processUI()
		},
	}

	cmdUI.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose logging. (default: false)")
	viper.BindPFlag("verbose", cmdUI.PersistentFlags().Lookup("verbose"))

	return cmdUI
}
