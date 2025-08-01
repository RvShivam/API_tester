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

var deleteCmd = &cobra.Command{
	Use:   "delete [URL]",
	Short: "Send a DELETE request to the specified URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		headers := parseHeaders(headersFlag)

		opts := internal.RequestOptions{
			Method:  "DELETE",
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

		// Format JSON response if possible
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
	deleteCmd.Flags().StringVar(&headersFlag, "headers", "", "Comma-separated headers (key:value,key:value)")
	deleteCmd.Flags().StringVar(&authFlag, "auth", "", "Authorization header (e.g., 'Bearer token' or 'Basic base64')")
	rootCmd.AddCommand(deleteCmd)
}
