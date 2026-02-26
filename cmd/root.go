package cmd

import (
	"fmt"
	"os"

	"github.com/RvShivam/API_tester/internal"
	"github.com/spf13/cobra"
)

var (
	envFile string
	Env     internal.Env
)

var rootCmd = &cobra.Command{
	Use:   "apitester",
	Short: "A CLI-based API testing tool",
	Long:  `A lightweight terminal-based API tester that supports REST methods, headers, body, authentication, and environment configs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		env, err := internal.LoadEnv(envFile)
		if err != nil {
			return err
		}
		Env = env
		if envFile != "" {
			fmt.Printf("Loaded environment: %s (%d variables)\n", envFile, len(env))
		}
		return nil
	},
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

func init() {
	rootCmd.PersistentFlags().StringVar(&envFile, "env", "", "Path to an environment JSON file (e.g., dev.json)")
}
