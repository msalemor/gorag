package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func verCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gorag\nVersion: %s\n", Version)
		},
	}
}
