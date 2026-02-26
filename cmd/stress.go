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
	stressConcurrencyFlag int
	stressDurationFlag    string
	stressRequestsFlag    int
	stressBodyFlag        string
	stressHeadersFlag     string
	stressAuthFlag        string
	stressMethodFlag      string
)

var stressCmd = &cobra.Command{
	Use:   "stress [URL]",
	Short: "Run a stress/load test against a URL",
	Long: `Hammer an API endpoint with concurrent requests to measure its performance.

Reports: total requests, successes, failures, requests/sec, and latency
percentiles (Min, Max, Avg, P50, P95, P99).`,
	Example: `  apitester stress https://httpbin.org/get --concurrency 20 --duration 15s
  apitester stress https://api.example.com/data --method GET --concurrency 10 --requests 500
  apitester stress "{{base_url}}/users" --env dev.json --concurrency 30 --duration 30s`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		url = Env.Interpolate(url)
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		method := strings.ToUpper(stressMethodFlag)

		body := stressBodyFlag
		body = Env.Interpolate(body)

		if body != "" {
			if err := internal.ValidateJSON(body); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}

		headers := parseHeaders(stressHeadersFlag)
		for k, v := range headers {
			headers[k] = Env.Interpolate(v)
		}

		// Parse duration
		duration, err := time.ParseDuration(stressDurationFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid duration %q: %v\n", stressDurationFlag, err)
			return
		}

		// If --requests is set, use a large duration and rely on the counter
		if stressRequestsFlag > 0 {
			duration = 24 * time.Hour // effectively unlimited time; workers stop via counter
		}

		opts := internal.StressOptions{
			Method:      method,
			URL:         url,
			Headers:     headers,
			Body:        body,
			Auth:        Env.Interpolate(stressAuthFlag),
			Concurrency: stressConcurrencyFlag,
			Duration:    duration,
			MaxRequests: stressRequestsFlag,
			Timeout:     10 * time.Second,
		}

		displayDuration := duration
		if stressRequestsFlag > 0 {
			displayDuration = 0 // will be overridden in report
		}

		fmt.Printf("🔥 Starting stress test → %s %s\n", method, url)
		fmt.Printf("   Concurrency: %d  |  ", opts.Concurrency)
		if stressRequestsFlag > 0 {
			fmt.Printf("Max Requests: %d\n", stressRequestsFlag)
		} else {
			fmt.Printf("Duration: %s\n", displayDuration)
		}
		fmt.Println()

		start := time.Now()
		result := internal.RunStress(opts)
		elapsed := time.Since(start)

		// Patch opts.Duration for the report if --requests was used
		if stressRequestsFlag > 0 {
			opts.Duration = elapsed
		}

		internal.PrintStressReport(opts, result)
	},
}

func init() {
	stressCmd.Flags().IntVar(&stressConcurrencyFlag, "concurrency", 10, "Number of concurrent workers")
	stressCmd.Flags().StringVar(&stressDurationFlag, "duration", "10s", "Duration of the test (e.g. 10s, 1m, 30s)")
	stressCmd.Flags().IntVar(&stressRequestsFlag, "requests", 0, "Total number of requests to send (overrides --duration)")
	stressCmd.Flags().StringVar(&stressMethodFlag, "method", "GET", "HTTP method to use")
	stressCmd.Flags().StringVar(&stressBodyFlag, "body", "", "JSON body for each request")
	stressCmd.Flags().StringVar(&stressHeadersFlag, "headers", "", "Comma-separated headers (key:value,...)")
	stressCmd.Flags().StringVar(&stressAuthFlag, "auth", "", "Auth header value")

	rootCmd.AddCommand(stressCmd)
}
