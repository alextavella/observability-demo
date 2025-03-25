package cmd

import (
	"fmt"
	"observability_demo/internal/frontend"

	"github.com/spf13/cobra"
)

var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "frontend is a cmd for managing frontend",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("frontend started")
		frontend.Start()
	},
}
