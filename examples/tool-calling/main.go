// Example: Tool calling with PALs
//
// This example demonstrates how to:
// - Create a custom tool with parameters
// - Create a PAL with the tool attached
// - List tools and their configurations
//
// Run with: TAVUS_API_KEY=your-key go run main.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/plexusone/tavus-go"
	"github.com/plexusone/tavus-go/internal/api"
)

func main() {
	client, err := tavus.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// First, get a face to use
	faces, err := client.Faces().List(ctx, api.ListFacesParams{Limit: api.NewOptInt(1)})
	if err != nil {
		log.Fatalf("Failed to list faces: %v", err)
	}
	if len(faces.Data) == 0 {
		log.Fatal("No faces available")
	}
	faceID := faces.Data[0].FaceID.Value

	// Create a tool for getting weather information
	fmt.Println("=== Creating Weather Tool ===")

	tool, err := client.Tools().Create(ctx, &api.CreateToolRequest{
		Name:        "get_weather",
		Description: "Get current weather for a location",
	})
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	fmt.Printf("Tool created: %s (%s)\n", tool.Name.Value, tool.ToolID.Value)

	// List all tools
	fmt.Println("\n=== Available Tools ===")
	tools, err := client.Tools().List(ctx, api.ListToolsParams{})
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, t := range tools.Data {
		fmt.Printf("  %s (%s): %s\n", t.Name.Value, t.ToolID.Value, t.Description.Value)
	}

	// Create a PAL that uses the tool
	fmt.Println("\n=== Creating PAL with Tool ===")
	pal, err := client.Pals().Create(ctx, &api.CreatePalRequest{
		PalName:       "Weather Assistant",
		DefaultFaceID: faceID,
		SystemPrompt: api.NewOptString(`You are a helpful weather assistant.
When users ask about weather, use the get_weather tool to fetch current conditions.
Be friendly and provide helpful weather-related advice.`),
	})
	if err != nil {
		log.Fatalf("Failed to create PAL: %v", err)
	}

	fmt.Printf("PAL created: %s\n", pal.PalID.Value)

	// Attach the tool to the PAL
	fmt.Println("\n=== Attaching Tool to PAL ===")
	if err := client.Pals().AttachTool(ctx, pal.PalID.Value, tool.ToolID.Value); err != nil {
		log.Fatalf("Failed to attach tool: %v", err)
	}
	fmt.Println("Tool attached successfully!")

	// Cleanup: Delete the PAL and tool
	fmt.Println("\n=== Cleanup ===")
	if err := client.Pals().Delete(ctx, pal.PalID.Value); err != nil {
		log.Printf("Warning: Failed to delete PAL: %v", err)
	} else {
		fmt.Println("PAL deleted")
	}

	if err := client.Tools().Delete(ctx, tool.ToolID.Value); err != nil {
		log.Printf("Warning: Failed to delete tool: %v", err)
	} else {
		fmt.Println("Tool deleted")
	}
}
