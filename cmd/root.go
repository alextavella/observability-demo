package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "observability-demo",
	Short: "A demo for a observability environment",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(frontendCmd)
	rootCmd.AddCommand(backendCmd)
	rootCmd.AddCommand(producerCmd)
	rootCmd.AddCommand(consumerCmd)
}
