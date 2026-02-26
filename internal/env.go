package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Env holds a map of key-value environment variables loaded from a file.
type Env map[string]string

// varPattern matches template variables like {{variable_name}}.
var varPattern = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// LoadEnv reads a JSON file and returns an Env map.
// Expected file format:
//
//	{
//	  "base_url": "https://api.example.com",
//	  "auth_token": "my-secret-token"
//	}
func LoadEnv(filename string) (Env, error) {
	if filename == "" {
		return Env{}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read env file %q: %w", filename, err)
	}

	var env Env
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("invalid JSON in env file %q: %w", filename, err)
	}

	return env, nil
}

// Interpolate replaces all {{variable}} placeholders in the input string
// with values from the Env map. If a variable is not found in the map,
// the placeholder is left as-is and a warning is printed.
func (e Env) Interpolate(input string) string {
	return varPattern.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name from {{name}}
		inner := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := e[inner]; ok {
			return val
		}
		fmt.Printf("Warning: environment variable %q not found, keeping placeholder\n", inner)
		return match
	})
}
