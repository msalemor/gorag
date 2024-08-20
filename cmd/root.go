package cmd

import (
	"github.com/spf13/cobra"
)

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
