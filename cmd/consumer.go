package cmd

import (
	"fmt"
	"observability_demo/internal/consumer"

	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "consumer is a cmd for managing consumer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consumer started")
		consumer.Start()
	},
}
