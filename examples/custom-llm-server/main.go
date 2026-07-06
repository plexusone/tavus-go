// Example: Custom LLM Server (OpenAI-compatible)
//
// This example demonstrates how to implement an OpenAI-compatible
// chat completions endpoint that Tavus can use as a custom LLM backend.
//
// The server:
// - Exposes POST /chat/completions endpoint
// - Validates API key authentication
// - Proxies requests to OpenAI (or your own LLM)
// - Supports streaming responses
//
// Usage:
//   1. Start this server: go run main.go
//   2. Expose via ngrok: ngrok http 8000
//   3. Configure your Tavus PAL to use the ngrok URL as the LLM base_url
//
// Based on: github.com/Tavus-Engineering/tavus-examples/cvi-custom-llm-with-backend
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ChatCompletionRequest represents the incoming request
type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionChunk represents a streaming response chunk
type ChatCompletionChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int          `json:"index"`
	Delta        Delta        `json:"delta"`
	FinishReason *string      `json:"finish_reason"`
}

// Delta represents the delta content in streaming
type Delta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/chat/completions", handleChatCompletions)
	http.HandleFunc("/health", handleHealth)

	log.Printf("Custom LLM server starting on port %s", port)
	log.Printf("Endpoint: POST /chat/completions")
	log.Printf("Health: GET /health")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log request info
	log.Printf("Request: %s %s", r.Method, r.URL.Path)
	log.Printf("Headers: %v", r.Header)

	// Authenticate request
	if !authenticateRequest(r) {
		log.Println("Authentication failed")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("Authentication successful")

	// Parse request body
	var req ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to parse request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Messages: %+v", req.Messages)

	// Check for OpenAI API key
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		// If no OpenAI key, return a mock response
		log.Println("No OPENAI_API_KEY, returning mock response")
		handleMockResponse(w, req)
		return
	}

	// Proxy to OpenAI
	handleOpenAIProxy(w, req, openaiKey)
}

func authenticateRequest(r *http.Request) bool {
	// Check Authorization header (Bearer token)
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "" {
			return true
		}
	}

	// Check X-API-Key header
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		return true
	}

	return false
}

func handleMockResponse(w http.ResponseWriter, req ChatCompletionRequest) {
	// Set headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Mock response content
	response := "Hello! I'm a mock LLM response. In production, this would be powered by your actual LLM backend. How can I help you today?"

	// Stream the response word by word
	words := strings.Fields(response)
	for i, word := range words {
		chunk := ChatCompletionChunk{
			ID:      "mock-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   "mock-model",
			Choices: []Choice{
				{
					Index: 0,
					Delta: Delta{Content: word + " "},
				},
			},
		}

		data, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()

		// Add slight delay for realistic streaming
		if i < len(words)-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	// Send done signal
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func handleOpenAIProxy(w http.ResponseWriter, req ChatCompletionRequest, apiKey string) {
	// Prepare OpenAI request
	req.Stream = true
	req.Model = "gpt-4o-mini" // Default model

	body, _ := json.Marshal(req)
	openaiReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		log.Printf("Failed to create OpenAI request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	openaiReq.Header.Set("Content-Type", "application/json")
	openaiReq.Header.Set("Authorization", "Bearer "+apiKey)

	// Make request to OpenAI
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(openaiReq)
	if err != nil {
		log.Printf("OpenAI request failed: %v", err)
		http.Error(w, "Upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("OpenAI error: %s", string(body))
		http.Error(w, "Upstream error", resp.StatusCode)
		return
	}

	// Set headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Stream response from OpenAI to client
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			fmt.Fprintf(w, "%s\n", line)
			flusher.Flush()
		}
	}
}
