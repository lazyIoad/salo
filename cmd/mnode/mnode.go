package main

import (
	"fmt"
	"os"

	"github.com/lazyIoad/salo/internal/mnode"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "salom",
	Short: "salo managed node server",
	Run: func(cmd *cobra.Command, args []string) {
		s := mnode.NewApiServer("/tmp/salo/server.sock")
		s.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to start worker server: %w", err))
		os.Exit(1)
	}
}
