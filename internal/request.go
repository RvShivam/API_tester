package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// RequestOptions holds all configurable parameters for an HTTP request.
type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Auth    string
	Timeout time.Duration
}

// SendRequest dispatches an HTTP request and returns the response, body bytes, duration, and any error.
// The duration includes the full round-trip including reading the response body.
func SendRequest(opts RequestOptions) (*http.Response, []byte, time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	var bodyReader io.Reader
	if opts.Body != "" {
		bodyReader = bytes.NewBufferString(opts.Body)
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, bodyReader)
	if err != nil {
		return nil, nil, 0, err
	}

	// Set default content type for requests with a body.
	// Check canonical header form to avoid duplicates.
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

	client := &http.Client{}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, time.Since(start), err
	}
	defer resp.Body.Close()

	// Read the full body BEFORE stopping the timer so duration is accurate.
	body, err := io.ReadAll(resp.Body)
	duration := time.Since(start)

	if err != nil {
		return resp, nil, duration, err
	}

	return resp, body, duration, nil
}

// ReadBodyInteractive reads a JSON body from stdin with a user prompt.
// Returns empty string if nothing is entered.
func ReadBodyInteractive() (string, error) {
	fmt.Println("Enter JSON body (end with Ctrl+D or Ctrl+Z on Windows):")
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(os.Stdin); err != nil {
		return "", fmt.Errorf("error reading request body: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}

// ValidateJSON returns an error if the provided string is not valid JSON.
func ValidateJSON(s string) error {
	var js interface{}
	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

// PrintResponse pretty-prints an HTTP response body to stdout.
func PrintResponse(resp *http.Response, body []byte, duration time.Duration) {
	fmt.Printf("Status:   %s\n", resp.Status)
	fmt.Printf("Duration: %v\n", duration)

	var pretty bytes.Buffer
	if json.Indent(&pretty, body, "", "  ") == nil {
		fmt.Println("Response (JSON):")
		fmt.Println(pretty.String())
	} else {
		fmt.Println("Response (raw):")
		fmt.Println(string(body))
	}
}
