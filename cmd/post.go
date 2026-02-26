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
	bodyFlag string
)

var postCmd = &cobra.Command{
	Use:   "post [URL]",
	Short: "Send a POST request to the specified URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
		url = Env.Interpolate(url)

		body := bodyFlag
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
			Method:  "POST",
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
	postCmd.Flags().StringVar(&headersFlag, "headers", "", "Comma-separated headers (key:value,key:value)")
	postCmd.Flags().StringVar(&authFlag, "auth", "", "Authorization header (e.g., 'Bearer token' or 'Basic base64')")
	postCmd.Flags().StringVar(&bodyFlag, "body", "", "JSON body for the request")
	rootCmd.AddCommand(postCmd)
}
