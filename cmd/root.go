package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apitester",
	Short: "A CLI-based API testing tool",
	Long:  `A lightweight terminal-based API tester that supports REST methods, headers, body, authentication, and environment configs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("API Tester CLI — Use 'apitester help' to get started.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
