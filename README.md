# Tavus-Go SDK

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/tavus-go/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/tavus-go/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/tavus-go/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/tavus-go/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/tavus-go/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/tavus-go/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/tavus-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/tavus-go
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/tavus-go
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Ftavus-go
 [loc-svg]: https://tokei.rs/b1/github/plexusone/tavus-go
 [repo-url]: https://github.com/plexusone/tavus-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/tavus-go/blob/main/LICENSE

Go SDK for the [Tavus API](https://docs.tavus.io/).

Tavus provides conversational video interfaces (CVI), video generation, and AI avatar management capabilities.

## Installation

```bash
go get github.com/plexusone/tavus-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/plexusone/tavus-go"
    "github.com/plexusone/tavus-go/api"
)

func main() {
    // Create client (reads TAVUS_API_KEY from environment if not provided)
    client, err := tavus.NewClient(tavus.WithAPIKey("your-api-key"))
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create a conversation
    conv, err := client.Conversations().Create(ctx, &api.CreateConversationRequest{
        PalID:  api.NewOptString("pal_abc123"),
        FaceID: api.NewOptString("face_xyz789"),
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Conversation created: %s", conv.ConversationID.Value)
    log.Printf("Join URL: %s", conv.ConversationURL.Value)
}
```

## Configuration

### Options

```go
// Use environment variable TAVUS_API_KEY
client, _ := tavus.NewClient()

// Explicit API key
client, _ := tavus.NewClient(tavus.WithAPIKey("your-api-key"))

// Custom base URL
client, _ := tavus.NewClient(tavus.WithBaseURL("https://custom.tavusapi.com"))

// Custom timeout
client, _ := tavus.NewClient(tavus.WithTimeout(60 * time.Second))

// Custom HTTP client
client, _ := tavus.NewClient(tavus.WithHTTPClient(&http.Client{
    Transport: customTransport,
}))
```

## Services

### Conversations

Real-time video conversation sessions.

```go
// Create a conversation
conv, err := client.Conversations().Create(ctx, &api.CreateConversationRequest{
    PalID:  api.NewOptString("pal_id"),
    FaceID: api.NewOptString("face_id"),
})

// List conversations
list, err := client.Conversations().List(ctx, api.ListConversationsParams{
    Limit: api.NewOptInt(10),
})

// Get a conversation
conv, err := client.Conversations().Get(ctx, "conversation_id")

// End a conversation
err := client.Conversations().End(ctx, "conversation_id")

// Delete a conversation
err := client.Conversations().Delete(ctx, "conversation_id")
```

### PALs (Personality AI Layers)

AI personas with custom instructions and behaviors.

```go
// Create a PAL
pal, err := client.Pals().Create(ctx, &api.CreatePalRequest{
    Name:               "Sales Assistant",
    SystemInstructions: api.NewOptString("You are a helpful sales assistant..."),
})

// List PALs
list, err := client.Pals().List(ctx, api.ListPalsParams{})

// Get a PAL
pal, err := client.Pals().Get(ctx, "pal_id")

// Update a PAL
pal, err := client.Pals().Update(ctx, "pal_id", &api.UpdatePalRequest{
    Name: api.NewOptString("Updated Name"),
})

// Attach a tool to a PAL
err := client.Pals().AttachTool(ctx, "pal_id", "tool_id")

// Delete a PAL
err := client.Pals().Delete(ctx, "pal_id")
```

### Faces

Visual identities for avatars.

```go
// Create/train a face
face, err := client.Faces().Create(ctx, &api.CreateFaceRequest{
    TrainingVideo: "https://example.com/training-video.mp4",
    FaceName:      api.NewOptString("My Avatar"),
})

// List faces
list, err := client.Faces().List(ctx, api.ListFacesParams{})

// Get a face
face, err := client.Faces().Get(ctx, "face_id")

// Rename a face
err := client.Faces().Rename(ctx, "face_id", "New Name")

// Delete a face
err := client.Faces().Delete(ctx, "face_id")
```

### Voices

Stock voices for avatars.

```go
// List available voices
voices, err := client.Voices().List(ctx, api.ListVoicesParams{})
```

### Videos

Async video generation.

```go
// Create a video
video, err := client.Videos().Create(ctx, &api.CreateVideoRequest{
    ReplicaID: "replica_id",
    Script:    "Hello, this is a generated video.",
})

// List videos
list, err := client.Videos().List(ctx, api.ListVideosParams{})

// Get a video
video, err := client.Videos().Get(ctx, "video_id")

// Rename a video
err := client.Videos().Rename(ctx, "video_id", "New Name")

// Delete a video
err := client.Videos().Delete(ctx, "video_id")
```

### Tools

Function calling tools for PALs.

```go
// Create a tool
tool, err := client.Tools().Create(ctx, &api.CreateToolRequest{
    Name:        "get_weather",
    Description: api.NewOptString("Get current weather for a location"),
    Parameters:  api.NewOptCreateToolRequestParameters(params),
})

// List tools
list, err := client.Tools().List(ctx, api.ListToolsParams{})

// Get a tool
tool, err := client.Tools().Get(ctx, "tool_id")

// Update a tool
tool, err := client.Tools().Update(ctx, "tool_id", &api.UpdateToolRequest{
    Description: api.NewOptString("Updated description"),
})

// Delete a tool
err := client.Tools().Delete(ctx, "tool_id")
```

### Guardrails

Behavioral boundaries for conversations.

```go
// Create a guardrail
guardrail, err := client.Guardrails().Create(ctx, &api.CreateGuardrailRequest{
    Name:        "No Personal Info",
    Description: api.NewOptString("Don't ask for personal information"),
})

// List guardrails
list, err := client.Guardrails().List(ctx, api.ListGuardrailsParams{})

// Get a guardrail
guardrail, err := client.Guardrails().Get(ctx, "guardrail_id")

// Update a guardrail
guardrail, err := client.Guardrails().Update(ctx, "guardrail_id", &api.UpdateGuardrailRequest{
    Name: api.NewOptString("Updated Name"),
})

// Delete a guardrail
err := client.Guardrails().Delete(ctx, "guardrail_id")
```

### Objectives

Conversation goals and outcomes.

```go
// Create objectives
objectives, err := client.Objectives().Create(ctx, &api.CreateObjectivesRequest{
    Objectives: []api.Objective{
        {Name: "Book Demo", Description: api.NewOptString("Schedule a demo call")},
    },
})

// List objectives
list, err := client.Objectives().List(ctx, api.ListObjectivesParams{})

// Get objectives
objectives, err := client.Objectives().Get(ctx, "objectives_id")

// Update objectives
objectives, err := client.Objectives().Update(ctx, "objectives_id", &api.UpdateObjectivesRequest{})

// Delete objectives
err := client.Objectives().Delete(ctx, "objectives_id")
```

### Documents

Knowledge base documents for PALs.

```go
// Create/upload a document
doc, err := client.Documents().Create(ctx, &api.CreateDocumentRequest{
    URL:      api.NewOptString("https://example.com/document.pdf"),
    FileName: api.NewOptString("product-guide.pdf"),
})

// List documents
list, err := client.Documents().List(ctx, api.ListDocumentsParams{})

// Get a document
doc, err := client.Documents().Get(ctx, "document_id")

// Update document metadata
doc, err := client.Documents().Update(ctx, "document_id", &api.UpdateDocumentRequest{
    FileName: api.NewOptString("updated-name.pdf"),
})

// Delete a document
err := client.Documents().Delete(ctx, "document_id")
```

### Deployments

Distribution channels for conversations.

```go
// Create a deployment
deployment, err := client.Deployments().Create(ctx, &api.CreateDeploymentRequest{
    Name:  "Production Widget",
    PalID: "pal_id",
})

// List deployments
list, err := client.Deployments().List(ctx, api.ListDeploymentsParams{})

// Get a deployment
deployment, err := client.Deployments().Get(ctx, "deployment_id")

// Update a deployment
deployment, err := client.Deployments().Update(ctx, "deployment_id", &api.UpdateDeploymentRequest{
    Name: api.NewOptString("Updated Name"),
})

// Delete a deployment
err := client.Deployments().Delete(ctx, "deployment_id")
```

## Low-Level API Access

For operations not covered by the high-level wrapper methods, use the underlying ogen-generated client:

```go
client, _ := tavus.NewClient()
apiClient := client.API()

// Direct access to all generated methods
res, err := apiClient.SomeOperation(ctx, params)
```

## Development

### Prerequisites

- Go 1.21+
- [ogen](https://github.com/ogen-go/ogen) for code generation
- [ogen-tools](https://github.com/plexusone/ogen-tools) for post-processing

### Regenerate API Client

```bash
./generate.sh
```

This regenerates the client from `openapi/openapi.yaml` and applies post-processing fixes.

## License

MIT
