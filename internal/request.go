package internal

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"
)

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Auth    string
	Timeout time.Duration
}

func SendRequest(opts RequestOptions) (*http.Response, []byte, time.Duration, error) {
	client := &http.Client{
		Timeout: opts.Timeout,
	}

	var bodyReader io.Reader
	if opts.Body != "" {
		bodyReader = bytes.NewBufferString(opts.Body)
	}

	req, err := http.NewRequest(opts.Method, opts.URL, bodyReader)
	if err != nil {
		return nil, nil, 0, err
	}

	// Set default content type for JSON requests
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
			// If no prefix, assume Bearer token
			req.Header.Set("Authorization", "Bearer "+opts.Auth)
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return nil, nil, duration, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, duration, err
	}

	return resp, body, duration, nil
}
