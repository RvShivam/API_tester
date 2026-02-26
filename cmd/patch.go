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

		body := patchBodyFlag
		if body == "" {
			var err error
			body, err = internal.ReadBodyInteractive()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}

		if body != "" {
			if err := internal.ValidateJSON(body); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}

		opts := internal.RequestOptions{
			Method:  "PATCH",
			URL:     url,
			Headers: parseHeaders(headersFlag),
			Body:    body,
			Auth:    authFlag,
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
	patchCmd.Flags().StringVar(&headersFlag, "headers", "", "Comma-separated headers (key:value,key:value)")
	patchCmd.Flags().StringVar(&authFlag, "auth", "", "Authorization header (e.g., 'Bearer token' or 'Basic base64')")
	patchCmd.Flags().StringVar(&patchBodyFlag, "body", "", "JSON body for the request")
	rootCmd.AddCommand(patchCmd)
}
