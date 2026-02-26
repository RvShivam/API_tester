package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/RvShivam/API_tester/internal"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [URL]",
	Short: "Send a DELETE request to the specified URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		url = Env.Interpolate(url)
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		headers := parseHeaders(headersFlag)
		for k, v := range headers {
			headers[k] = Env.Interpolate(v)
		}

		opts := internal.RequestOptions{
			Method:  "DELETE",
			URL:     url,
			Headers: headers,
			Auth:    Env.Interpolate(authFlag),
			Timeout: 10 * time.Second,
		}

		resp, body, duration, err := internal.SendRequest(opts)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		internal.PrintResponse(resp, body, duration)
	},
}

func init() {
	deleteCmd.Flags().StringVar(&headersFlag, "headers", "", "Comma-separated headers (key:value,key:value)")
	deleteCmd.Flags().StringVar(&authFlag, "auth", "", "Authorization header (e.g., 'Bearer token' or 'Basic base64')")
	rootCmd.AddCommand(deleteCmd)
}
