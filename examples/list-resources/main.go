// Example: List all Tavus resources
//
// This example demonstrates how to list all resource types
// available in your Tavus account.
//
// Run with: TAVUS_API_KEY=your-key go run main.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/plexusone/tavus-go"
	"github.com/plexusone/tavus-go/api"
)

func main() {
	client, err := tavus.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Voices (stock voices)
	fmt.Println("=== Voices ===")
	voices, err := client.Voices().List(ctx, api.ListVoicesParams{})
	if err != nil {
		log.Printf("Failed to list voices: %v", err)
	} else {
		fmt.Printf("Found %d voices\n", len(voices.Data))
		for _, v := range voices.Data {
			fmt.Printf("  - %s\n", v.VoiceName.Value)
		}
	}

	// Faces
	fmt.Println("\n=== Faces ===")
	faces, err := client.Faces().List(ctx, api.ListFacesParams{Limit: api.NewOptInt(10)})
	if err != nil {
		log.Printf("Failed to list faces: %v", err)
	} else {
		fmt.Printf("Found %d faces (showing first 10)\n", len(faces.Data))
		for _, f := range faces.Data {
			fmt.Printf("  - %s (%s) [%s]\n", f.FaceName.Value, f.FaceID.Value, f.Status.Value)
		}
	}

	// PALs
	fmt.Println("\n=== PALs ===")
	pals, err := client.Pals().List(ctx, api.ListPalsParams{Limit: api.NewOptInt(10)})
	if err != nil {
		log.Printf("Failed to list PALs: %v", err)
	} else {
		fmt.Printf("Found %d PALs (showing first 10)\n", len(pals.Data))
		for _, p := range pals.Data {
			fmt.Printf("  - %s (%s)\n", p.PalName.Value, p.PalID.Value)
		}
	}

	// Conversations
	fmt.Println("\n=== Recent Conversations ===")
	convs, err := client.Conversations().List(ctx, api.ListConversationsParams{Limit: api.NewOptInt(10)})
	if err != nil {
		log.Printf("Failed to list conversations: %v", err)
	} else {
		fmt.Printf("Found %d conversations (showing first 10)\n", len(convs.Data))
		for _, c := range convs.Data {
			fmt.Printf("  - %s [%s]\n", c.ConversationID.Value, c.Status.Value)
		}
	}

	// Tools
	fmt.Println("\n=== Tools ===")
	tools, err := client.Tools().List(ctx, api.ListToolsParams{})
	if err != nil {
		log.Printf("Failed to list tools: %v", err)
	} else {
		fmt.Printf("Found %d tools\n", len(tools.Data))
		for _, t := range tools.Data {
			fmt.Printf("  - %s (%s)\n", t.Name.Value, t.ToolID.Value)
		}
	}

	// Guardrails
	fmt.Println("\n=== Guardrails ===")
	guardrails, err := client.Guardrails().List(ctx, api.ListGuardrailsParams{})
	if err != nil {
		log.Printf("Failed to list guardrails: %v", err)
	} else {
		fmt.Printf("Found %d guardrails\n", len(guardrails.Data))
		for _, g := range guardrails.Data {
			fmt.Printf("  - %s (%s)\n", g.GuardrailName.Value, g.UUID.Value)
		}
	}

	// Objectives
	fmt.Println("\n=== Objectives ===")
	objectives, err := client.Objectives().List(ctx, api.ListObjectivesParams{})
	if err != nil {
		log.Printf("Failed to list objectives: %v", err)
	} else {
		fmt.Printf("Found %d objective sets\n", len(objectives.Data))
	}

	// Documents
	fmt.Println("\n=== Documents ===")
	docs, err := client.Documents().List(ctx, api.ListDocumentsParams{})
	if err != nil {
		log.Printf("Failed to list documents: %v", err)
	} else {
		fmt.Printf("Found %d documents\n", len(docs.Data))
		for _, d := range docs.Data {
			fmt.Printf("  - %s (%s)\n", d.DocumentName.Value, d.DocumentID.Value)
		}
	}

	// Deployments
	fmt.Println("\n=== Deployments ===")
	deployments, err := client.Deployments().List(ctx, api.ListDeploymentsParams{})
	if err != nil {
		log.Printf("Failed to list deployments: %v", err)
	} else {
		fmt.Printf("Found %d deployments\n", len(deployments.Data))
		for _, d := range deployments.Data {
			fmt.Printf("  - %s (%s)\n", d.Name.Value, d.DeploymentID.Value)
		}
	}

	// Videos
	fmt.Println("\n=== Videos ===")
	videos, err := client.Videos().List(ctx, api.ListVideosParams{Limit: api.NewOptInt(10)})
	if err != nil {
		log.Printf("Failed to list videos: %v", err)
	} else {
		fmt.Printf("Found %d videos (showing first 10)\n", len(videos.Data))
		for _, v := range videos.Data {
			fmt.Printf("  - %s (%s) [%s]\n", v.VideoName.Value, v.VideoID.Value, v.Status.Value)
		}
	}
}
