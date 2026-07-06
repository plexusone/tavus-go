# tavus-go Examples

This directory contains example applications demonstrating how to use the tavus-go SDK.

## Prerequisites

Set your Tavus API key:

```bash
export TAVUS_API_KEY=your-api-key
```

## Examples

### basic

Basic conversation creation and management. Demonstrates:

- Creating a Tavus client
- Listing available faces and PALs
- Creating and ending a conversation

```bash
cd basic
go run main.go
```

### list-resources

List all resources in your Tavus account. Useful for exploring what's available.

```bash
cd list-resources
go run main.go
```

### pal-management

PAL (Personality AI Layer) management. Demonstrates:

- Creating a PAL with custom instructions
- Creating guardrails for behavioral boundaries
- Creating objectives for conversation goals
- Updating PAL settings

```bash
cd pal-management
go run main.go
```

### tool-calling

Function calling with PALs. Demonstrates:

- Creating tools with JSON schema parameters
- Attaching tools to PALs
- Listing and managing tools

```bash
cd tool-calling
go run main.go
```

## LiveKit Integration

For LiveKit transport integration examples, see the [omni-livekit](https://github.com/plexusone/omni-livekit) project which demonstrates real-time avatar sessions using LiveKit.
