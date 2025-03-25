package cmd

import (
	"fmt"
	"observability_demo/internal/backend"

	"github.com/spf13/cobra"
)

var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "backend is a cmd for managing backend",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("backend started")
		backend.Start()
	},
}
