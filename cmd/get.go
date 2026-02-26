package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/RvShivam/API_tester/internal"
	"github.com/spf13/cobra"
)

var (
	headersFlag string
	authFlag    string
)

var getCmd = &cobra.Command{
	Use:   "get [URL]",
	Short: "Send a GET request to a URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		opts := internal.RequestOptions{
			Method:  "GET",
			URL:     url,
			Headers: parseHeaders(headersFlag),
			Auth:    authFlag,
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
	getCmd.Flags().StringVar(&headersFlag, "headers", "", "Comma-separated headers (key:value,key:value)")
	getCmd.Flags().StringVar(&authFlag, "auth", "", "Authorization header (e.g., 'Bearer token' or 'Basic base64')")
	rootCmd.AddCommand(getCmd)
}

func parseHeaders(input string) map[string]string {
	headers := make(map[string]string)
	if input == "" {
		return headers
	}
	pairs := strings.Split(input, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return headers
}
