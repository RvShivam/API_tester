# Stage 1: Build the Go web server
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire repo root
COPY . .

# Build the apitester binary (the main CLI tool)
RUN go build -o apitester .

# Build the web server
WORKDIR /app/web
RUN go mod download
RUN go build -o server .

# Stage 2: Create the final lean image
FROM alpine:latest

WORKDIR /app

# Copy the apitester CLI binary from the builder
COPY --from=builder /app/apitester ./apitester

# Copy the web server binary
COPY --from=builder /app/web/server ./server

# Copy the static frontend files
COPY --from=builder /app/web/public ./public

# Copy the demo environment file used by the --env demo button
COPY --from=builder /app/web/demo-env.json ./demo-env.json

RUN chmod +x ./apitester ./server

EXPOSE 8080

CMD ["./server"]
