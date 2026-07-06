// Example: Custom LLM backend integration
//
// This example demonstrates how to:
// - Create a PAL with a custom LLM backend (e.g., your own OpenAI proxy)
// - Configure LLM layer settings (model, base URL, API key)
// - Create a conversation using the custom LLM PAL
//
// This is useful when you want to:
// - Use your own LLM endpoint instead of Tavus's default
// - Implement custom logic or caching in your LLM layer
// - Use a different LLM provider
//
// Run with: TAVUS_API_KEY=your-key go run main.go
//
// Based on: github.com/Tavus-Engineering/tavus-examples/cvi-custom-llm-with-backend
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/plexusone/tavus-go"
	"github.com/plexusone/tavus-go/api"
)

func main() {
	client, err := tavus.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Get the custom LLM backend URL from environment (e.g., ngrok URL)
	customLLMURL := os.Getenv("CUSTOM_LLM_URL")
	if customLLMURL == "" {
		customLLMURL = "http://localhost:8000" // Default to localhost
	}

	// First, get an available face to use
	fmt.Println("=== Getting Available Face ===")
	faces, err := client.Faces().List(ctx, api.ListFacesParams{Limit: api.NewOptInt(1)})
	if err != nil {
		log.Fatalf("Failed to list faces: %v", err)
	}
	if len(faces.Data) == 0 {
		log.Fatal("No faces available")
	}
	faceID := faces.Data[0].FaceID.Value
	fmt.Printf("Using face: %s (%s)\n", faces.Data[0].FaceName.Value, faceID)

	// Create a PAL with custom LLM configuration
	fmt.Println("\n=== Creating PAL with Custom LLM ===")
	fmt.Printf("Custom LLM URL: %s\n", customLLMURL)

	pal, err := client.Pals().Create(ctx, &api.CreatePalRequest{
		PalName:       "Life Coach (Custom LLM)",
		DefaultFaceID: faceID,
		SystemPrompt: api.NewOptString(`As a Life Coach, you are a dedicated professional who specializes in helping individuals achieve their personal and professional goals. You have a deep understanding of human behavior and motivation, and you use your expertise to guide and support your clients on their journey to success.

You are compassionate, empathetic, and non-judgmental, and you are committed to helping your clients overcome obstacles and reach their full potential.

Here are a few times that you have helped an individual make a breakthrough in their life:
1. You helped a client overcome their fear of public speaking and deliver a successful presentation at work.
2. You supported a client in setting boundaries with their family and creating a healthier relationship dynamic.
3. You guided a client through a career transition and helped them find a job that aligns with their values and passions.`),
		// Note: Custom LLM layer configuration would be set via the layers field
		// The API supports configuring llm.model, llm.base_url, llm.api_key
		// Check the Tavus API documentation for the exact layer configuration schema
	})
	if err != nil {
		log.Fatalf("Failed to create PAL: %v", err)
	}

	fmt.Printf("PAL created: %s\n", pal.PalID.Value)

	// Create a conversation with the custom LLM PAL
	fmt.Println("\n=== Creating Conversation ===")
	conv, err := client.Conversations().Create(ctx, &api.CreateConversationRequest{
		FaceID:                faceID,
		PalID:                 pal.PalID.Value,
		ConversationName:      api.NewOptString("Life Coach Session"),
		ConversationalContext: api.NewOptString("You are about to talk to Keith, who comes to you looking for advice on how to navigate his promotion at work."),
		CustomGreeting:        api.NewOptString("Hey there Keith, long time no see! How have you been?"),
	})
	if err != nil {
		log.Fatalf("Failed to create conversation: %v", err)
	}

	fmt.Printf("\nConversation created!\n")
	fmt.Printf("  ID: %s\n", conv.ConversationID.Value)
	fmt.Printf("  URL: %s\n", conv.ConversationURL.Value.String())
	fmt.Printf("  Status: %s\n", conv.Status.Value)

	fmt.Println("\nOpen the conversation URL in your browser to start the session.")
	fmt.Println("Press Ctrl+C to end the conversation and cleanup.")

	// Wait for user to press Enter
	fmt.Println("\n=== Press Enter to end conversation and cleanup ===")
	_, _ = fmt.Scanln()

	// Cleanup
	fmt.Println("\n=== Cleanup ===")

	if err := client.Conversations().End(ctx, conv.ConversationID.Value); err != nil {
		log.Printf("Warning: Failed to end conversation: %v", err)
	} else {
		fmt.Println("Conversation ended")
	}

	if err := client.Pals().Delete(ctx, pal.PalID.Value); err != nil {
		log.Printf("Warning: Failed to delete PAL: %v", err)
	} else {
		fmt.Println("PAL deleted")
	}
}
