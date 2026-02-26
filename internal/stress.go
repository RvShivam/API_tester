package internal

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// StressOptions defines the parameters for a stress test run.
type StressOptions struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        string
	Auth        string
	Concurrency int
	Duration    time.Duration
	MaxRequests int // 0 means unlimited (use Duration instead)
	Timeout     time.Duration
}

// StressResult holds the aggregated results of a stress test.
type StressResult struct {
	TotalRequests int
	Successes     int
	Failures      int
	Latencies     []time.Duration
	Errors        []string
}

// rateLimitedClient does a single HTTP request via a reusable client.
var stressClient = &http.Client{}

// RunStress executes a load test against a URL using a goroutine worker pool.
func RunStress(opts StressOptions) StressResult {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Duration+5*time.Second)
	defer cancel()

	type result struct {
		latency time.Duration
		err     error
		status  int
	}

	resultCh := make(chan result, opts.Concurrency*10)
	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Stop after Duration
	go func() {
		time.Sleep(opts.Duration)
		close(stop)
	}()

	// Shared request counter for MaxRequests mode
	var (
		mu      sync.Mutex
		counter int
	)

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				return
			default:
			}

			if opts.MaxRequests > 0 {
				mu.Lock()
				if counter >= opts.MaxRequests {
					mu.Unlock()
					return
				}
				counter++
				mu.Unlock()
			}

			req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, strings.NewReader(opts.Body))
			if err != nil {
				resultCh <- result{err: err}
				continue
			}

			if opts.Body != "" && req.Header.Get("Content-Type") == "" {
				req.Header.Set("Content-Type", "application/json")
			}
			for k, v := range opts.Headers {
				req.Header.Set(k, v)
			}
			if opts.Auth != "" {
				if strings.HasPrefix(opts.Auth, "Bearer ") || strings.HasPrefix(opts.Auth, "Basic ") {
					req.Header.Set("Authorization", opts.Auth)
				} else {
					req.Header.Set("Authorization", "Bearer "+opts.Auth)
				}
			}

			start := time.Now()
			resp, err := stressClient.Do(req)
			if err != nil {
				resultCh <- result{latency: time.Since(start), err: err}
				continue
			}
			latency := time.Since(start)
			resp.Body.Close()

			resultCh <- result{latency: latency, status: resp.StatusCode}
		}
	}

	for i := 0; i < opts.Concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	// Close the result channel once all workers finish
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var sr StressResult
	for r := range resultCh {
		sr.TotalRequests++
		if r.err != nil {
			sr.Failures++
			errMsg := r.err.Error()
			if len(sr.Errors) < 5 { // store only first 5 unique errors
				sr.Errors = append(sr.Errors, errMsg)
			}
		} else if r.status >= 200 && r.status < 400 {
			sr.Successes++
			sr.Latencies = append(sr.Latencies, r.latency)
		} else {
			sr.Failures++
		}
	}

	return sr
}

// PrintStressReport prints a formatted summary report to stdout.
func PrintStressReport(opts StressOptions, result StressResult) {
	fmt.Println()
	fmt.Println("════════════════════ STRESS TEST REPORT ════════════════════")
	fmt.Printf("  Target:       %s %s\n", opts.Method, opts.URL)
	fmt.Printf("  Concurrency:  %d workers\n", opts.Concurrency)
	fmt.Printf("  Duration:     %s\n", opts.Duration)
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Printf("  Total Reqs:   %d\n", result.TotalRequests)
	fmt.Printf("  Successes:    %d\n", result.Successes)
	fmt.Printf("  Failures:     %d\n", result.Failures)
	if result.TotalRequests > 0 && opts.Duration > 0 {
		rps := float64(result.TotalRequests) / opts.Duration.Seconds()
		fmt.Printf("  Req/sec:      %.2f\n", rps)
	}

	if len(result.Latencies) > 0 {
		sort.Slice(result.Latencies, func(i, j int) bool {
			return result.Latencies[i] < result.Latencies[j]
		})
		n := len(result.Latencies)
		var total time.Duration
		for _, l := range result.Latencies {
			total += l
		}
		avg := total / time.Duration(n)

		fmt.Println("────────────────────────────────────────────────────────────")
		fmt.Printf("  Latency Min:  %v\n", result.Latencies[0])
		fmt.Printf("  Latency Max:  %v\n", result.Latencies[n-1])
		fmt.Printf("  Latency Avg:  %v\n", avg)
		fmt.Printf("  Latency P50:  %v\n", result.Latencies[percentileIdx(n, 50)])
		fmt.Printf("  Latency P95:  %v\n", result.Latencies[percentileIdx(n, 95)])
		fmt.Printf("  Latency P99:  %v\n", result.Latencies[percentileIdx(n, 99)])
	}

	if len(result.Errors) > 0 {
		fmt.Println("────────────────────────────────────────────────────────────")
		fmt.Println("  Sample Errors:")
		for _, e := range result.Errors {
			fmt.Printf("    • %s\n", e)
		}
	}
	fmt.Println("════════════════════════════════════════════════════════════")
}

func percentileIdx(n int, p float64) int {
	idx := int(math.Ceil(float64(n)*p/100)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}
	return idx
}
