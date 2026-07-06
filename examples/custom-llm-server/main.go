// Example: Custom LLM Server (OpenAI-compatible)
//
// This example demonstrates how to implement an OpenAI-compatible
// chat completions endpoint that Tavus can use as a custom LLM backend.
//
// The server:
// - Exposes POST /chat/completions endpoint
// - Validates API key authentication
// - Uses omnillm-core for multi-provider LLM support (OpenAI, Anthropic, etc.)
// - Supports streaming responses
//
// Environment variables:
// - PORT: Server port (default: 8000)
// - LLM_PROVIDER: Provider name (openai, anthropic, xai, ollama, etc.)
// - LLM_MODEL: Model to use (e.g., gpt-4o-mini, claude-3-5-haiku)
// - OPENAI_API_KEY: OpenAI API key (when using openai provider)
// - ANTHROPIC_API_KEY: Anthropic API key (when using anthropic provider)
// - XAI_API_KEY: X.AI API key (when using xai provider)
// - OLLAMA_BASE_URL: Ollama base URL (when using ollama provider)
//
// Usage:
//  1. Start this server: LLM_PROVIDER=openai OPENAI_API_KEY=... go run main.go
//  2. Expose via ngrok: ngrok http 8000
//  3. Configure your Tavus PAL to use the ngrok URL as the LLM base_url
//
// Based on: github.com/Tavus-Engineering/tavus-examples/cvi-custom-llm-with-backend
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	omnillm "github.com/plexusone/omnillm-core"
	"github.com/plexusone/omnillm-core/provider"
)

// Global LLM client
var llmClient *omnillm.ChatClient

// Config holds server configuration
type Config struct {
	Port     string
	Provider omnillm.ProviderName
	Model    string
	APIKey   string
	BaseURL  string
}

// ChatCompletionChunk represents a streaming response chunk (OpenAI format)
type ChatCompletionChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

// Delta represents the delta content in streaming
type Delta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

func main() {
	config := loadConfig()

	// Initialize LLM client
	var err error
	llmClient, err = createLLMClient(config)
	if err != nil {
		log.Printf("Warning: Failed to create LLM client: %v", err)
		log.Println("Server will return mock responses")
	} else {
		log.Printf("LLM client initialized: provider=%s, model=%s", config.Provider, config.Model)
	}

	http.HandleFunc("/chat/completions", handleChatCompletions)
	http.HandleFunc("/health", handleHealth)

	log.Printf("Custom LLM server starting on port %s", config.Port)
	log.Printf("Endpoint: POST /chat/completions")
	log.Printf("Health: GET /health")

	server := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func loadConfig() Config {
	config := Config{
		Port:     getEnvDefault("PORT", "8000"),
		Provider: omnillm.ProviderName(getEnvDefault("LLM_PROVIDER", "openai")),
		Model:    getEnvDefault("LLM_MODEL", ""),
	}

	// Get API key and base URL based on provider
	switch config.Provider {
	case omnillm.ProviderNameOpenAI:
		config.APIKey = os.Getenv("OPENAI_API_KEY")
		if config.Model == "" {
			config.Model = omnillm.ModelGPT4oMini
		}
	case omnillm.ProviderNameAnthropic:
		config.APIKey = os.Getenv("ANTHROPIC_API_KEY")
		if config.Model == "" {
			config.Model = omnillm.ModelClaude3_5Haiku
		}
	case omnillm.ProviderNameXAI:
		config.APIKey = os.Getenv("XAI_API_KEY")
		if config.Model == "" {
			config.Model = omnillm.ModelGrok3Mini
		}
	case omnillm.ProviderNameOllama:
		config.BaseURL = getEnvDefault("OLLAMA_BASE_URL", "http://localhost:11434")
		if config.Model == "" {
			config.Model = omnillm.ModelOllamaLlama3_8B
		}
	case omnillm.ProviderNameKimi:
		config.APIKey = os.Getenv("KIMI_API_KEY")
		if config.Model == "" {
			config.Model = omnillm.ModelMoonshotV1_8K
		}
	case omnillm.ProviderNameGLM:
		config.APIKey = os.Getenv("GLM_API_KEY")
		if config.Model == "" {
			config.Model = omnillm.ModelGLM4_5Flash
		}
	case omnillm.ProviderNameQwen:
		config.APIKey = os.Getenv("QWEN_API_KEY")
		if config.Model == "" {
			config.Model = omnillm.ModelQwenFlash
		}
	}

	return config
}

func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func createLLMClient(config Config) (*omnillm.ChatClient, error) {
	// Ollama doesn't require an API key
	if config.Provider != omnillm.ProviderNameOllama && config.APIKey == "" {
		return nil, fmt.Errorf("no API key for provider %s", config.Provider)
	}

	return omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{
				Provider: config.Provider,
				APIKey:   config.APIKey,
				BaseURL:  config.BaseURL,
			},
		},
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log request info
	log.Printf("Request: %s %s", r.Method, r.URL.Path) //nolint:gosec // G706: URL path is safe to log

	// Authenticate request
	if !authenticateRequest(r) {
		log.Println("Authentication failed")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req provider.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to parse request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Model: %s, Messages: %d", req.Model, len(req.Messages))

	// Use mock response if no LLM client
	if llmClient == nil {
		handleMockResponse(w, req)
		return
	}

	// Handle streaming request via omnillm
	handleOmnillmStream(w, req)
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
	return apiKey != ""
}

func handleMockResponse(w http.ResponseWriter, req provider.ChatCompletionRequest) {
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
	response := "Hello! I'm a mock LLM response. Configure LLM_PROVIDER and the appropriate API key to use a real LLM backend."

	// Use the requested model or fallback to mock-model
	model := req.Model
	if model == "" {
		model = "mock-model"
	}

	// Stream the response word by word
	words := strings.Fields(response)
	for i, word := range words {
		chunk := ChatCompletionChunk{
			ID:      fmt.Sprintf("mock-%d", time.Now().UnixNano()),
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   model,
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

func handleOmnillmStream(w http.ResponseWriter, req provider.ChatCompletionRequest) {
	ctx := context.Background()

	// Create streaming request
	stream, err := llmClient.CreateChatCompletionStream(ctx, &req)
	if err != nil {
		log.Printf("Failed to create stream: %v", err)
		http.Error(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	// Set headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Stream response chunks
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Stream error: %v", err)
			break
		}

		// Convert to OpenAI-compatible format
		openaiChunk := ChatCompletionChunk{
			ID:      chunk.ID,
			Object:  "chat.completion.chunk",
			Created: chunk.Created,
			Model:   chunk.Model,
			Choices: make([]Choice, len(chunk.Choices)),
		}

		for i, c := range chunk.Choices {
			openaiChunk.Choices[i] = Choice{
				Index:        c.Index,
				FinishReason: c.FinishReason,
			}
			if c.Delta != nil {
				openaiChunk.Choices[i].Delta = Delta{
					Content: c.Delta.Content,
					Role:    string(c.Delta.Role),
				}
			}
		}

		data, _ := json.Marshal(openaiChunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	// Send done signal
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
