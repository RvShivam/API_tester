package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RvShivam/API_tester/internal"
	"github.com/spf13/cobra"
)

var (
	putBodyFlag string
)

var putCmd = &cobra.Command{
	Use:   "put [URL]",
	Short: "Send a PUT request to the specified URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
		url = Env.Interpolate(url)

		body := putBodyFlag
		if body == "" {
			var err error
			body, err = internal.ReadBodyInteractive()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		body = Env.Interpolate(body)

		if body != "" {
			if err := internal.ValidateJSON(body); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}

		headers := parseHeaders(headersFlag)
		for k, v := range headers {
			headers[k] = Env.Interpolate(v)
		}

		opts := internal.RequestOptions{
			Method:  "PUT",
			URL:     url,
			Headers: headers,
			Body:    body,
			Auth:    Env.Interpolate(authFlag),
			Timeout: 15 * time.Second,
		}

		resp, respBody, duration, err := internal.SendRequest(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			return
		}

		internal.PrintResponse(resp, respBody, duration)
	},
}

func init() {
	putCmd.Flags().StringVar(&headersFlag, "headers", "", "Comma-separated headers (key:value,key:value)")
	putCmd.Flags().StringVar(&authFlag, "auth", "", "Authorization header (e.g., 'Bearer token' or 'Basic base64')")
	putCmd.Flags().StringVar(&putBodyFlag, "body", "", "JSON body for the request")
	rootCmd.AddCommand(putCmd)
}
