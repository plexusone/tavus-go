// Example: PAL (Personality AI Layer) management
//
// This example demonstrates how to:
// - Create a PAL with custom instructions
// - Configure PAL layers (LLM, TTS, STT)
// - Update PAL settings
// - Add guardrails and objectives
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

	// Create a PAL for a sales assistant
	fmt.Println("=== Creating Sales Assistant PAL ===")
	pal, err := client.Pals().Create(ctx, &api.CreatePalRequest{
		PalName:       "Sales Demo Assistant",
		DefaultFaceID: faceID,
		SystemPrompt: api.NewOptString(`You are a friendly and professional sales assistant for Acme Corp.
Your goals are to:
1. Understand the customer's needs
2. Explain how our products can help them
3. Answer questions about pricing and features
4. Schedule demo calls when appropriate

Be conversational but professional. Ask clarifying questions.
Never make up information about products or pricing.`),
	})
	if err != nil {
		log.Fatalf("Failed to create PAL: %v", err)
	}

	fmt.Printf("PAL created: %s (%s)\n", pal.PalID.Value, "Sales Demo Assistant")

	// Get full PAL details
	fmt.Println("\n=== PAL Details ===")
	palDetails, err := client.Pals().Get(ctx, pal.PalID.Value)
	if err != nil {
		log.Fatalf("Failed to get PAL: %v", err)
	}
	fmt.Printf("  Name: %s\n", palDetails.PalName.Value)
	fmt.Printf("  ID: %s\n", palDetails.PalID.Value)
	fmt.Printf("  Pipeline Mode: %s\n", palDetails.PipelineMode.Value)

	// Create a guardrail for the PAL
	fmt.Println("\n=== Creating Guardrail ===")
	guardrail, err := client.Guardrails().Create(ctx, &api.CreateGuardrailRequest{
		GuardrailName: "No Competitor Mentions",
		GuardrailPrompt: `Do not discuss or compare products from competitors.
If asked about competitors, politely redirect the conversation to our products.
Say: "I'm here to help you learn about Acme Corp products. Let me tell you about what we offer."`,
	})
	if err != nil {
		log.Fatalf("Failed to create guardrail: %v", err)
	}
	fmt.Printf("Guardrail created: %s\n", guardrail.UUID.Value)

	// Create objectives for the conversation
	fmt.Println("\n=== Creating Objectives ===")
	objectives, err := client.Objectives().Create(ctx, &api.CreateObjectivesRequest{
		Objectives: []api.ObjectiveItem{
			{
				ObjectiveName:   "qualify_lead",
				ObjectivePrompt: "Understand the customer's company size, budget, and timeline",
			},
			{
				ObjectiveName:   "product_fit",
				ObjectivePrompt: "Determine which product tier best fits their needs",
			},
			{
				ObjectiveName:   "schedule_demo",
				ObjectivePrompt: "Book a follow-up demo call with the sales team",
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create objectives: %v", err)
	}
	fmt.Printf("Objectives created: %s\n", objectives.ObjectivesID.Value)

	// Update the PAL with a new name
	fmt.Println("\n=== Updating PAL ===")
	updatedPal, err := client.Pals().Update(ctx, pal.PalID.Value, &api.UpdatePalRequest{
		PalName: api.NewOptString("Acme Sales Assistant v2"),
	})
	if err != nil {
		log.Fatalf("Failed to update PAL: %v", err)
	}
	fmt.Printf("PAL updated: %s\n", updatedPal.PalName.Value)

	// List all PALs
	fmt.Println("\n=== All PALs ===")
	pals, err := client.Pals().List(ctx, api.ListPalsParams{Limit: api.NewOptInt(10)})
	if err != nil {
		log.Fatalf("Failed to list PALs: %v", err)
	}
	for _, p := range pals.Data {
		fmt.Printf("  %s (%s)\n", p.PalName.Value, p.PalID.Value)
	}

	// Cleanup
	fmt.Println("\n=== Cleanup ===")

	if err := client.Objectives().Delete(ctx, objectives.ObjectivesID.Value); err != nil {
		log.Printf("Warning: Failed to delete objectives: %v", err)
	} else {
		fmt.Println("Objectives deleted")
	}

	if err := client.Guardrails().Delete(ctx, guardrail.UUID.Value); err != nil {
		log.Printf("Warning: Failed to delete guardrail: %v", err)
	} else {
		fmt.Println("Guardrail deleted")
	}

	if err := client.Pals().Delete(ctx, pal.PalID.Value); err != nil {
		log.Printf("Warning: Failed to delete PAL: %v", err)
	} else {
		fmt.Println("PAL deleted")
	}
}
