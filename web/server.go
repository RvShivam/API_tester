package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for the demo. Tighten this in production.
		return true
	},
}

// allowedCommands defines the ONLY sub-commands the web terminal is allowed to execute.
// This prevents shell injection attacks.
var allowedCommands = map[string]bool{
	"get":        true,
	"post":       true,
	"put":        true,
	"delete":     true,
	"patch":      true,
	"stress":     true,
	"collection": true,
	"help":       true,
	"--help":     true,
	"-h":         true,
	"version":    true,
	"--version":  true,
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// parseArgs splits a command string into tokens the same way a shell would,
// stripping surrounding single or double quotes from each token.
// e.g.: `get https://x.com --headers "Accept:application/json"` →
//        ["get", "https://x.com", "--headers", "Accept:application/json"]
func parseArgs(s string) []string {
	var args []string
	var cur strings.Builder
	inQuote := false
	var quoteChar rune

	for _, r := range s {
		switch {
		case inQuote:
			if r == quoteChar {
				inQuote = false
			} else {
				cur.WriteRune(r)
			}
		case r == '"' || r == '\'':
			inQuote = true
			quoteChar = r
		case r == ' ' || r == '\t':
			if cur.Len() > 0 {
				args = append(args, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		args = append(args, cur.String())
	}
	return args
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()
	log.Println("New WebSocket connection established")

	// Send a welcome message to the client on connect.
	welcome := "\r\n\033[1;32m Welcome to API Tester CLI — Web Terminal \033[0m\r\n" +
		"\033[90m Type a command below or click a demo to get started.\033[0m\r\n\r\n"
	conn.WriteMessage(websocket.TextMessage, []byte(welcome))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// Ignore normal disconnect errors (1001 Going Away, 1000 Normal Closure, 1006 Abnormal Closure).
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		rawCmd := strings.TrimSpace(string(msg))
		if rawCmd == "" {
			continue
		}

		log.Printf("Received command: %s", rawCmd)

		// Echo the command back so it displays in the terminal.
		conn.WriteMessage(websocket.TextMessage, []byte("\r\n\033[1;36m$ apitester "+rawCmd+"\033[0m\r\n"))

		// Security: parse args and strip leading 'apitester' prefix if typed.
		// This allows both 'get https://...' and 'apitester get https://...' to work.
		parts := parseArgs(rawCmd)
		if len(parts) > 0 && (parts[0] == "apitester" || parts[0] == "apitester.exe") {
			parts = parts[1:]
		}
		if len(parts) == 0 {
			conn.WriteMessage(websocket.TextMessage, []byte("\033[90mHint: try 'get https://httpbin.org/get' or '--help'\033[0m\r\n"))
			continue
		}
		if !allowedCommands[parts[0]] {
			errMsg := fmt.Sprintf("\033[1;31m[BLOCKED] '%s' is not a recognised apitester command.\033[0m\r\n", parts[0])
			conn.WriteMessage(websocket.TextMessage, []byte(errMsg))
			continue
		}

		// Resolve the path to the apitester binary.
		// Check current dir first, then parent dir (repo root).
		// Also check .exe variants for Windows compatibility.
		candidates := []string{"./apitester", "./apitester.exe", "../apitester", "../apitester.exe"}
		binaryPath := ""
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				binaryPath = c
				break
			}
		}
		if binaryPath == "" {
			conn.WriteMessage(websocket.TextMessage, []byte("\033[1;31m[ERROR] apitester binary not found.\033[0m\r\n"))
			continue
		}

		// Execute the apitester binary with a 60-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		cmd := exec.CommandContext(ctx, binaryPath, parts...)
		output, err := cmd.CombinedOutput()

		// Convert newlines for the terminal (\n -> \r\n).
		result := strings.ReplaceAll(string(output), "\n", "\r\n")

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				conn.WriteMessage(websocket.TextMessage, []byte("\033[1;31m[ERROR] Command timed out after 60s.\033[0m\r\n"))
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte(result))
			}
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte(result))
		}

		conn.WriteMessage(websocket.TextMessage, []byte("\r\n"))
		cancel() // release context resources after each command
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve the static frontend files from the ./public directory.
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
