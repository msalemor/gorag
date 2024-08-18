package cmd

import (
	"fmt"
	"log"

	"github.com/msalemor/gorag/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func consoleCmd() *cobra.Command {
	cmdConsole := &cobra.Command{
		Use:   "console",
		Short: "Run in console mode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running in console mode")
			fmt.Println("Config:", viper.GetString("config"))
			log.Println("Verbose:", Verbose)
			process.ProcessConsole(OllamaChatEndpoint, OllamaEmbeddingEndpoint, DEFAULT_COLLECTION, EmbeddingModel, Keep, Verbose)
		},
	}
	cmdConsole.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose logging. (default: false)")
	viper.BindPFlag("verbose", cmdConsole.PersistentFlags().Lookup("verbose"))

	cmdConsole.PersistentFlags().BoolVarP(&Keep, "keep", "k", false, "Keep the collection. (default: false)")
	viper.BindPFlag("keep", cmdConsole.PersistentFlags().Lookup("keep"))

	return cmdConsole
}
