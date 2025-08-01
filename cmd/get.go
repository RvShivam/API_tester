package cmd

import (
	"bytes"
	"encoding/json"
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

		headers := parseHeaders(headersFlag)

		opts := internal.RequestOptions{
			Method:  "GET",
			URL:     url,
			Headers: headers,
			Auth:    authFlag,
			Timeout: 10 * time.Second,
		}

		resp, body, duration, err := internal.SendRequest(opts)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Status:", resp.Status)
		fmt.Println("Time taken:", duration)

		var pretty bytes.Buffer
		if json.Indent(&pretty, body, "", "  ") == nil {
			fmt.Println("Response (JSON):")
			fmt.Println(pretty.String())
		} else {
			fmt.Println("Response (raw):")
			fmt.Println(string(body))
		}
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
