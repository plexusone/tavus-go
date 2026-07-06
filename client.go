// Package tavus provides a Go client for the Tavus API.
//
// Tavus is an AI video research company providing conversational video interfaces,
// video generation, and avatar management capabilities.
//
// The client wraps the ogen-generated API client with a higher-level interface
// that handles authentication and provides convenient methods for common operations.
//
// # Quick Start
//
//	client, err := tavus.NewClient(tavus.WithAPIKey("your-api-key"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create a conversation
//	conv, err := client.Conversations().Create(ctx, &api.CreateConversationRequest{
//	    PalId:  "pal-id",
//	    FaceId: "face-id",
//	})
//
// # Environment Variables
//
// If no API key is provided, the client will look for TAVUS_API_KEY in the environment.
package tavus

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/plexusone/tavus-go/internal/api"
)

// Version is the SDK version.
const Version = "0.1.0"

// DefaultBaseURL is the default Tavus API base URL.
const DefaultBaseURL = "https://tavusapi.com"

// Client is the main Tavus client for interacting with the API.
type Client struct {
	apiClient *api.Client
	apiKey    string
	baseURL   string

	// Domain-based service accessors
	conversationsSvc *ConversationsService
	palsSvc          *PalsService
	facesSvc         *FacesService
	voicesSvc        *VoicesService
	videosSvc        *VideosService
	toolsSvc         *ToolsService
	guardrailsSvc    *GuardrailsService
	objectivesSvc    *ObjectivesService
	documentsSvc     *DocumentsService
	deploymentsSvc   *DeploymentsService
}

// NewClient creates a new Tavus client with the given options.
func NewClient(opts ...Option) (*Client, error) {
	options := defaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Try environment variable if API key not set
	if options.apiKey == "" {
		options.apiKey = os.Getenv("TAVUS_API_KEY")
	}

	// Create HTTP client
	httpClient := options.httpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: options.timeout,
		}
	}

	// Create security source for authentication
	securitySource := &apiKeySecuritySource{apiKey: options.apiKey}

	// Create the ogen client
	apiClient, err := api.NewClient(
		options.baseURL,
		securitySource,
		api.WithClient(&sdkHTTPClient{client: httpClient}),
	)
	if err != nil {
		return nil, err
	}

	c := &Client{
		apiClient: apiClient,
		apiKey:    options.apiKey,
		baseURL:   options.baseURL,
	}

	// Initialize domain-based services
	c.conversationsSvc = &ConversationsService{client: apiClient}
	c.palsSvc = &PalsService{client: apiClient}
	c.facesSvc = &FacesService{client: apiClient}
	c.voicesSvc = &VoicesService{client: apiClient}
	c.videosSvc = &VideosService{client: apiClient}
	c.toolsSvc = &ToolsService{client: apiClient}
	c.guardrailsSvc = &GuardrailsService{client: apiClient}
	c.objectivesSvc = &ObjectivesService{client: apiClient}
	c.documentsSvc = &DocumentsService{client: apiClient}
	c.deploymentsSvc = &DeploymentsService{client: apiClient}

	return c, nil
}

// apiKeySecuritySource implements api.SecuritySource for API key authentication.
type apiKeySecuritySource struct {
	apiKey string
}

// ApiKeyAuth implements api.SecuritySource.
func (s *apiKeySecuritySource) ApiKeyAuth(ctx context.Context, operationName api.OperationName) (api.ApiKeyAuth, error) {
	return api.ApiKeyAuth{APIKey: s.apiKey}, nil
}

// sdkHTTPClient wraps an http.Client to add SDK version headers.
type sdkHTTPClient struct {
	client *http.Client
}

// Do implements ht.Client interface.
func (c *sdkHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Add SDK version headers
	req.Header.Set("X-Tavus-SDK-Version", Version)
	req.Header.Set("X-Tavus-SDK-Lang", "go")

	return c.client.Do(req)
}

// API returns the underlying ogen-generated API client for advanced usage.
// Use this when you need access to API endpoints not covered by the
// high-level wrapper methods.
func (c *Client) API() *api.Client {
	return c.apiClient
}

// Conversations returns the conversations service for real-time video sessions.
func (c *Client) Conversations() *ConversationsService {
	return c.conversationsSvc
}

// Pals returns the PALs service for personality AI layer management.
func (c *Client) Pals() *PalsService {
	return c.palsSvc
}

// Faces returns the faces service for visual identity management.
func (c *Client) Faces() *FacesService {
	return c.facesSvc
}

// Voices returns the voices service for stock voice listing.
func (c *Client) Voices() *VoicesService {
	return c.voicesSvc
}

// Videos returns the videos service for async video generation.
func (c *Client) Videos() *VideosService {
	return c.videosSvc
}

// Tools returns the tools service for function calling management.
func (c *Client) Tools() *ToolsService {
	return c.toolsSvc
}

// Guardrails returns the guardrails service for behavioral boundaries.
func (c *Client) Guardrails() *GuardrailsService {
	return c.guardrailsSvc
}

// Objectives returns the objectives service for conversation goals.
func (c *Client) Objectives() *ObjectivesService {
	return c.objectivesSvc
}

// Documents returns the documents service for knowledge base management.
func (c *Client) Documents() *DocumentsService {
	return c.documentsSvc
}

// Deployments returns the deployments service for distribution channels.
func (c *Client) Deployments() *DeploymentsService {
	return c.deploymentsSvc
}

// APIKey returns the API key used by the client.
func (c *Client) APIKey() string {
	return c.apiKey
}

// BaseURL returns the base URL used by the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// clientOptions holds the options for creating a Client.
type clientOptions struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		baseURL: DefaultBaseURL,
		timeout: 120 * time.Second,
	}
}

// Option is a functional option for configuring the Client.
type Option func(*clientOptions)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(apiKey string) Option {
	return func(o *clientOptions) {
		o.apiKey = apiKey
	}
}

// WithBaseURL sets the API base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}
