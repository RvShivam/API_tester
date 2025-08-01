# API Tester CLI 

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
![Go Version](https://img.shields.io/badge/Go-1.22-blue)
![Last Commit](https://img.shields.io/github/last-commit/RvShivam/API_tester)
![Top Lang](https://img.shields.io/github/languages/top/RvShivam/API_tester)


🚀 **API Tester** is a lightweight, terminal-based API testing tool built with Go and Cobra that emulates core features of Postman in a streamlined command-line interface. Perfect for developers who prefer working in the terminal or need to integrate API testing into scripts and automation workflows.


## ✨ Features

- **Full REST Method Support**: GET, POST, PUT, DELETE, PATCH requests
- **Flag-based Configuration**: Easy-to-use command-line flags for headers, body, and authentication
- **Smart JSON Handling**: Automatic JSON formatting and validation for request bodies
- **Pretty Response Formatting**: Automatically formats JSON responses with proper indentation
- **Request Timing**: Measures and displays request duration for performance analysis
- **Flexible Authentication**: Supports Bearer tokens, Basic auth, and custom authentication headers
- **Smart URL Handling**: Automatically adds HTTPS protocol if not specified
- **Cross-platform**: Works on Windows, macOS, and Linux

## 🔧 Installation

### Prerequisites
- Go 1.24.0 or later

### Build from Source

1. Clone or download the repository
2. Navigate to the project directory:
   ```sh
   cd path/to/API_tester
   ```
3. Build the executable:
   ```sh
   go build -o apitester.exe
   ```
   
   On Unix-like systems (macOS/Linux):
   ```sh
   go build -o apitester
   ```

## 📖 Usage

### Getting Help

To see all available commands:
```sh
apitester.exe --help
```

To get help for a specific command:
```sh
apitester.exe [command] --help
```

## 🔨 Available Commands

### GET Request
```sh
apitester.exe get [URL] [flags]
```

**Flags:**
- `--headers`: Comma-separated headers (key:value,key:value)
- `--auth`: Authorization header (e.g., 'Bearer token' or 'Basic base64')

### POST Request
```sh
apitester.exe post [URL] [flags]
```

**Flags:**
- `--body`: JSON body for the request
- `--headers`: Comma-separated headers (key:value,key:value)
- `--auth`: Authorization header (e.g., 'Bearer token' or 'Basic base64')

### PUT Request
```sh
apitester.exe put [URL] [flags]
```

**Flags:**
- `--body`: JSON body for the request
- `--headers`: Comma-separated headers (key:value,key:value)
- `--auth`: Authorization header (e.g., 'Bearer token' or 'Basic base64')

### DELETE Request
```sh
apitester.exe delete [URL] [flags]
```

**Flags:**
- `--headers`: Comma-separated headers (key:value,key:value)
- `--auth`: Authorization header (e.g., 'Bearer token' or 'Basic base64')

### PATCH Request
```sh
apitester.exe patch [URL] [flags]
```

**Flags:**
- `--body`: JSON body for the request
- `--headers`: Comma-separated headers (key:value,key:value)
- `--auth`: Authorization header (e.g., 'Bearer token' or 'Basic base64')

## 💡 Examples

### Simple GET Request
```sh
apitester.exe get https://httpbin.org/get
```

### GET Request with Headers and Authentication
```sh
apitester.exe get https://api.example.com/users --headers "Accept:application/json,User-Agent:MyApp/1.0" --auth "Bearer your_token_here"
```

### POST Request with JSON Body
```sh
apitester.exe post https://httpbin.org/post --body '{"name":"John Doe","email":"john@example.com"}' --headers "Content-Type:application/json"
```

### PUT Request with Authentication
```sh
apitester.exe put https://api.example.com/users/123 --body '{"name":"Updated Name"}' --auth "Bearer your_token_here"
```

### DELETE Request
```sh
apitester.exe delete https://api.example.com/users/123 --auth "Bearer your_token_here"
```

### PATCH Request
```sh
apitester.exe patch https://api.example.com/users/123 --body '{"status":"active"}' --headers "Content-Type:application/json"
```

## 🎯 Advanced Usage

### Multiple Headers
You can specify multiple headers separated by commas:
```sh
apitester.exe get https://api.example.com/data --headers "Accept:application/json,Authorization:Bearer token,X-Custom-Header:value"
```

### Authentication Types

**Bearer Token:**
```sh
--auth "Bearer your_jwt_token"
```

**Basic Authentication:**
```sh
--auth "Basic base64_encoded_credentials"
```

**Custom Authentication:**
If you don't specify "Bearer" or "Basic", the tool automatically assumes Bearer token:
```sh
--auth "your_token"  # Automatically becomes "Bearer your_token"
```

### Interactive Body Input
If you don't specify the `--body` flag for POST, PUT, or PATCH requests, the tool will prompt you to enter the JSON body interactively:
```sh
apitester.exe post https://httpbin.org/post
# Will prompt: "Enter JSON body (end with Ctrl+D or Ctrl+Z on Windows):"
```

## 📊 Response Format

The tool provides detailed response information including:

- **Status Code**: HTTP response status (e.g., "200 OK")
- **Request Duration**: Time taken to complete the request
- **Response Body**: Automatically formatted JSON (when applicable) or raw response

### Example Output
```
Status: 200 OK
Duration: 1.8596219s
Response (JSON):
{
  "args": {},
  "headers": {
    "Accept-Encoding": "gzip",
    "Host": "httpbin.org",
    "User-Agent": "Go-http-client/2.0"
  },
  "origin": "152.59.91.33",
  "url": "https://httpbin.org/get"
}
```

## 🛠️ Project Structure

```
API_tester/
├── cmd/
│   ├── root.go      # Root command definition
│   ├── get.go       # GET command implementation
│   ├── post.go      # POST command implementation
│   ├── put.go       # PUT command implementation
│   ├── delete.go    # DELETE command implementation
│   └── patch.go     # PATCH command implementation
├── internal/
│   └── request.go   # Core HTTP request handling logic
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
└── README.md        # This file
```

## License

This project is licensed under the [MIT License](LICENSE).


⭐ **Star this repository if you find it helpful!**
