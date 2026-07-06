package tavus

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("with API key option", func(t *testing.T) {
		client, err := NewClient(WithAPIKey("test-key"))
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if client.APIKey() != "test-key" {
			t.Errorf("expected API key 'test-key', got '%s'", client.APIKey())
		}
	})

	t.Run("with custom base URL", func(t *testing.T) {
		client, err := NewClient(
			WithAPIKey("test-key"),
			WithBaseURL("https://custom.tavusapi.com"),
		)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if client.BaseURL() != "https://custom.tavusapi.com" {
			t.Errorf("expected base URL 'https://custom.tavusapi.com', got '%s'", client.BaseURL())
		}
	})

	t.Run("with custom timeout", func(t *testing.T) {
		client, err := NewClient(
			WithAPIKey("test-key"),
			WithTimeout(30*time.Second),
		)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if client == nil {
			t.Error("expected non-nil client")
		}
	})

	t.Run("with custom HTTP client", func(t *testing.T) {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		client, err := NewClient(
			WithAPIKey("test-key"),
			WithHTTPClient(httpClient),
		)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if client == nil {
			t.Error("expected non-nil client")
		}
	})

	t.Run("reads from environment variable", func(t *testing.T) {
		t.Setenv("TAVUS_API_KEY", "env-test-key")

		client, err := NewClient()
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if client.APIKey() != "env-test-key" {
			t.Errorf("expected API key 'env-test-key', got '%s'", client.APIKey())
		}
	})

	t.Run("default base URL", func(t *testing.T) {
		client, err := NewClient(WithAPIKey("test-key"))
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if client.BaseURL() != DefaultBaseURL {
			t.Errorf("expected default base URL '%s', got '%s'", DefaultBaseURL, client.BaseURL())
		}
	})
}

func TestClientServices(t *testing.T) {
	client, err := NewClient(WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	t.Run("Conversations service", func(t *testing.T) {
		if client.Conversations() == nil {
			t.Error("expected non-nil ConversationsService")
		}
	})

	t.Run("Pals service", func(t *testing.T) {
		if client.Pals() == nil {
			t.Error("expected non-nil PalsService")
		}
	})

	t.Run("Faces service", func(t *testing.T) {
		if client.Faces() == nil {
			t.Error("expected non-nil FacesService")
		}
	})

	t.Run("Voices service", func(t *testing.T) {
		if client.Voices() == nil {
			t.Error("expected non-nil VoicesService")
		}
	})

	t.Run("Videos service", func(t *testing.T) {
		if client.Videos() == nil {
			t.Error("expected non-nil VideosService")
		}
	})

	t.Run("Tools service", func(t *testing.T) {
		if client.Tools() == nil {
			t.Error("expected non-nil ToolsService")
		}
	})

	t.Run("Guardrails service", func(t *testing.T) {
		if client.Guardrails() == nil {
			t.Error("expected non-nil GuardrailsService")
		}
	})

	t.Run("Objectives service", func(t *testing.T) {
		if client.Objectives() == nil {
			t.Error("expected non-nil ObjectivesService")
		}
	})

	t.Run("Documents service", func(t *testing.T) {
		if client.Documents() == nil {
			t.Error("expected non-nil DocumentsService")
		}
	})

	t.Run("Deployments service", func(t *testing.T) {
		if client.Deployments() == nil {
			t.Error("expected non-nil DeploymentsService")
		}
	})

	t.Run("API client", func(t *testing.T) {
		if client.API() == nil {
			t.Error("expected non-nil API client")
		}
	})
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}
