//go:build integration

package tavus

import (
	"context"
	"os"
	"testing"

	"github.com/plexusone/tavus-go/api"
)

// skipIfNoAPIKey skips the test if TAVUS_API_KEY is not set.
func skipIfNoAPIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("TAVUS_API_KEY") == "" {
		t.Skip("TAVUS_API_KEY not set, skipping integration test")
	}
}

// newTestClient creates a client for integration testing.
func newTestClient(t *testing.T) *Client {
	t.Helper()
	skipIfNoAPIKey(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestIntegration_Voices_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	voices, err := client.Voices().List(ctx, api.ListVoicesParams{})
	if err != nil {
		t.Fatalf("Voices.List failed: %v", err)
	}

	t.Logf("Found %d voices", len(voices.Data))
	for _, v := range voices.Data {
		t.Logf("  Voice: %s", v.VoiceName.Value)
	}
}

func TestIntegration_Faces_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	faces, err := client.Faces().List(ctx, api.ListFacesParams{})
	if err != nil {
		t.Fatalf("Faces.List failed: %v", err)
	}

	t.Logf("Found %d faces", len(faces.Data))
	for _, f := range faces.Data {
		t.Logf("  Face: %s (%s) - Status: %s", f.FaceName.Value, f.FaceID.Value, f.Status.Value)
	}
}

func TestIntegration_Pals_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	pals, err := client.Pals().List(ctx, api.ListPalsParams{})
	if err != nil {
		t.Fatalf("Pals.List failed: %v", err)
	}

	t.Logf("Found %d PALs", len(pals.Data))
	for _, p := range pals.Data {
		t.Logf("  PAL: %s (%s)", p.PalName.Value, p.PalID.Value)
	}
}

func TestIntegration_Conversations_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	convs, err := client.Conversations().List(ctx, api.ListConversationsParams{
		Limit: api.NewOptInt(10),
	})
	if err != nil {
		t.Fatalf("Conversations.List failed: %v", err)
	}

	t.Logf("Found %d conversations", len(convs.Data))
	for _, c := range convs.Data {
		t.Logf("  Conversation: %s - Status: %s", c.ConversationID.Value, c.Status.Value)
	}
}

func TestIntegration_Tools_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	tools, err := client.Tools().List(ctx, api.ListToolsParams{})
	if err != nil {
		t.Fatalf("Tools.List failed: %v", err)
	}

	t.Logf("Found %d tools", len(tools.Data))
	for _, tool := range tools.Data {
		t.Logf("  Tool: %s (%s)", tool.Name.Value, tool.ToolID.Value)
	}
}

func TestIntegration_Guardrails_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	guardrails, err := client.Guardrails().List(ctx, api.ListGuardrailsParams{})
	if err != nil {
		t.Fatalf("Guardrails.List failed: %v", err)
	}

	t.Logf("Found %d guardrails", len(guardrails.Data))
	for _, g := range guardrails.Data {
		t.Logf("  Guardrail: %s (%s)", g.GuardrailName.Value, g.UUID.Value)
	}
}

func TestIntegration_Objectives_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	objectives, err := client.Objectives().List(ctx, api.ListObjectivesParams{})
	if err != nil {
		t.Fatalf("Objectives.List failed: %v", err)
	}

	t.Logf("Found %d objectives", len(objectives.Data))
}

func TestIntegration_Documents_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	docs, err := client.Documents().List(ctx, api.ListDocumentsParams{})
	if err != nil {
		t.Fatalf("Documents.List failed: %v", err)
	}

	t.Logf("Found %d documents", len(docs.Data))
	for _, d := range docs.Data {
		t.Logf("  Document: %s (%s)", d.DocumentName.Value, d.DocumentID.Value)
	}
}

func TestIntegration_Deployments_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	deployments, err := client.Deployments().List(ctx, api.ListDeploymentsParams{})
	if err != nil {
		t.Fatalf("Deployments.List failed: %v", err)
	}

	t.Logf("Found %d deployments", len(deployments.Data))
	for _, d := range deployments.Data {
		t.Logf("  Deployment: %s (%s)", d.Name.Value, d.DeploymentID.Value)
	}
}

func TestIntegration_Videos_List(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	videos, err := client.Videos().List(ctx, api.ListVideosParams{})
	if err != nil {
		t.Fatalf("Videos.List failed: %v", err)
	}

	t.Logf("Found %d videos", len(videos.Data))
	for _, v := range videos.Data {
		t.Logf("  Video: %s (%s) - Status: %s", v.VideoName.Value, v.VideoID.Value, v.Status.Value)
	}
}
