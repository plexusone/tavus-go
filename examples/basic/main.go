// Example: Basic conversation creation and management
//
// This example demonstrates how to:
// - Create a Tavus client
// - List available faces and PALs
// - Create a conversation
// - Get conversation details
//
// Run with: TAVUS_API_KEY=your-key go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/plexusone/tavus-go"
	"github.com/plexusone/tavus-go/internal/api"
)

func main() {
	// Create client - reads TAVUS_API_KEY from environment
	client, err := tavus.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// List available faces
	fmt.Println("=== Available Faces ===")
	faces, err := client.Faces().List(ctx, api.ListFacesParams{
		Limit: api.NewOptInt(5),
	})
	if err != nil {
		log.Fatalf("Failed to list faces: %v", err)
	}
	for _, face := range faces.Data {
		fmt.Printf("  %s (%s) - %s\n", face.FaceName.Value, face.FaceID.Value, face.Status.Value)
	}

	// List available PALs
	fmt.Println("\n=== Available PALs ===")
	pals, err := client.Pals().List(ctx, api.ListPalsParams{
		Limit: api.NewOptInt(5),
	})
	if err != nil {
		log.Fatalf("Failed to list PALs: %v", err)
	}
	for _, pal := range pals.Data {
		fmt.Printf("  %s (%s)\n", pal.PalName.Value, pal.PalID.Value)
	}

	// Check if we have resources to create a conversation
	if len(faces.Data) == 0 || len(pals.Data) == 0 {
		fmt.Println("\nNo faces or PALs available. Create some first!")
		os.Exit(0)
	}

	// Create a conversation using the first available face and PAL
	faceID := faces.Data[0].FaceID.Value
	palID := pals.Data[0].PalID.Value

	fmt.Printf("\n=== Creating Conversation ===\n")
	fmt.Printf("Using Face: %s\n", faceID)
	fmt.Printf("Using PAL: %s\n", palID)

	conv, err := client.Conversations().Create(ctx, &api.CreateConversationRequest{
		FaceID: faceID,
		PalID:  palID,
	})
	if err != nil {
		log.Fatalf("Failed to create conversation: %v", err)
	}

	fmt.Printf("\nConversation created!\n")
	fmt.Printf("  ID: %s\n", conv.ConversationID.Value)
	fmt.Printf("  URL: %s\n", conv.ConversationURL.Value.String())
	fmt.Printf("  Status: %s\n", conv.Status.Value)

	// End the conversation (cleanup)
	fmt.Println("\n=== Ending Conversation ===")
	if err := client.Conversations().End(ctx, conv.ConversationID.Value); err != nil {
		log.Printf("Warning: Failed to end conversation: %v", err)
	} else {
		fmt.Println("Conversation ended successfully")
	}
}
