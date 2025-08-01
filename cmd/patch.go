package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RvShivam/API_tester/internal"
	"github.com/spf13/cobra"
)

var (
	patchBodyFlag string
)

var patchCmd = &cobra.Command{
	Use:   "patch [URL]",
	Short: "Send a PATCH request to the specified URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		var body string
		if patchBodyFlag != "" {
			body = patchBodyFlag
		} else {
			// Read the body from standard input
			fmt.Println("Enter JSON body (end with Ctrl+D or Ctrl+Z on Windows):")
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(os.Stdin); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
				return
			}
			body = strings.TrimSpace(buf.String())
		}

		// Validate JSON if provided
		if body != "" {
			var js map[string]interface{}
			if err := json.Unmarshal([]byte(body), &js); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid JSON body: %v\n", err)
				return
			}
		}

		headers := parseHeaders(headersFlag)

		opts := internal.RequestOptions{
			Method:  "PATCH",
			URL:     url,
			Headers: headers,
			Body:    body,
			Auth:    authFlag,
			Timeout: 15 * time.Second,
		}

		resp, respBody, duration, err := internal.SendRequest(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			return
		}

		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Duration: %v\n", duration)

		// Format JSON response if possible
		var pretty bytes.Buffer
		if json.Indent(&pretty, respBody, "", "  ") == nil {
			fmt.Println("Response (JSON):")
			fmt.Println(pretty.String())
		} else {
			fmt.Println("Response (raw):")
			fmt.Println(string(respBody))
		}
	},
}

func init() {
	patchCmd.Flags().StringVar(&headersFlag, "headers", "", "Comma-separated headers (key:value,key:value)")
	patchCmd.Flags().StringVar(&authFlag, "auth", "", "Authorization header (e.g., 'Bearer token' or 'Basic base64')")
	patchCmd.Flags().StringVar(&patchBodyFlag, "body", "", "JSON body for the request")
	rootCmd.AddCommand(patchCmd)
}
