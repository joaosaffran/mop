package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mob",
	Short: "A CLI tool to manage git workflow",
}

func Execute() error {
	return rootCmd.Execute()
}
