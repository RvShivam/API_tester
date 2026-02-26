package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RvShivam/API_tester/internal"
	"github.com/spf13/cobra"
)

// ── collection root ───────────────────────────────────────────────────────────

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Manage saved request collections",
	Long:  `Save, list, run, and delete named HTTP requests in your personal collection.`,
}

// ── collection save ───────────────────────────────────────────────────────────

var (
	saveNameFlag    string
	saveMethodFlag  string
	saveURLFlag     string
	saveHeadersFlag string
	saveBodyFlag    string
	saveAuthFlag    string
)

var collectionSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a request to the collection",
	Example: `  apitester collection save --name login --method POST \
    --url "{{base_url}}/auth/login" \
    --body '{"email":"user@example.com","password":"secret"}' \
    --auth "{{auth_token}}"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if saveNameFlag == "" {
			return fmt.Errorf("--name is required")
		}
		if saveMethodFlag == "" {
			return fmt.Errorf("--method is required")
		}
		if saveURLFlag == "" {
			return fmt.Errorf("--url is required")
		}

		method := strings.ToUpper(saveMethodFlag)

		req := internal.SavedRequest{
			Name:    saveNameFlag,
			Method:  method,
			URL:     saveURLFlag,
			Headers: parseHeaders(saveHeadersFlag),
			Body:    saveBodyFlag,
			Auth:    saveAuthFlag,
			Timeout: 15 * time.Second,
		}

		return internal.SaveRequest(req)
	},
}

// ── collection list ───────────────────────────────────────────────────────────

var collectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		requests, err := internal.ListRequests()
		if err != nil {
			return err
		}

		if len(requests) == 0 {
			fmt.Println("No saved requests. Use 'apitester collection save' to add one.")
			return nil
		}

		fmt.Printf("%-20s  %-7s  %s\n", "NAME", "METHOD", "URL")
		fmt.Println(strings.Repeat("─", 70))
		for _, r := range requests {
			fmt.Printf("%-20s  %-7s  %s\n", r.Name, r.Method, r.URL)
		}
		return nil
	},
}

// ── collection run ────────────────────────────────────────────────────────────

var collectionRunCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "Run a saved request by name",
	Example: `  apitester collection run login
  apitester collection run login --env dev.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		req, err := internal.GetRequest(name)
		if err != nil {
			return err
		}

		// Apply environment interpolation to all fields
		url := Env.Interpolate(req.URL)
		body := Env.Interpolate(req.Body)
		auth := Env.Interpolate(req.Auth)
		headers := make(map[string]string)
		for k, v := range req.Headers {
			headers[k] = Env.Interpolate(v)
		}

		timeout := req.Timeout
		if timeout == 0 {
			timeout = 15 * time.Second
		}

		fmt.Printf("Running %q [%s %s]\n\n", name, req.Method, url)

		opts := internal.RequestOptions{
			Method:  req.Method,
			URL:     url,
			Headers: headers,
			Body:    body,
			Auth:    auth,
			Timeout: timeout,
		}

		resp, respBody, duration, err := internal.SendRequest(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
			return nil
		}

		internal.PrintResponse(resp, respBody, duration)
		return nil
	},
}

// ── collection delete ─────────────────────────────────────────────────────────

var collectionDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a saved request by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := internal.DeleteRequest(name); err != nil {
			return err
		}
		fmt.Printf("🗑️  Deleted request %q from collection.\n", name)
		return nil
	},
}

// ── init ──────────────────────────────────────────────────────────────────────

func init() {
	// save flags
	collectionSaveCmd.Flags().StringVar(&saveNameFlag, "name", "", "Unique name for the request (required)")
	collectionSaveCmd.Flags().StringVar(&saveMethodFlag, "method", "", "HTTP method: GET, POST, PUT, DELETE, PATCH (required)")
	collectionSaveCmd.Flags().StringVar(&saveURLFlag, "url", "", "Request URL, supports {{variable}} syntax (required)")
	collectionSaveCmd.Flags().StringVar(&saveHeadersFlag, "headers", "", "Comma-separated headers (key:value,...)")
	collectionSaveCmd.Flags().StringVar(&saveBodyFlag, "body", "", "JSON body for the request")
	collectionSaveCmd.Flags().StringVar(&saveAuthFlag, "auth", "", "Auth header value")

	// register sub-commands
	collectionCmd.AddCommand(collectionSaveCmd)
	collectionCmd.AddCommand(collectionListCmd)
	collectionCmd.AddCommand(collectionRunCmd)
	collectionCmd.AddCommand(collectionDeleteCmd)

	// register with root
	rootCmd.AddCommand(collectionCmd)
}
