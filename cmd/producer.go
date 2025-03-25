package cmd

import (
	"fmt"
	"observability_demo/internal/producer"

	"github.com/spf13/cobra"
)

var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "producer is a cmd for managing producer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("producer started")
		producer.Start()
	},
}
