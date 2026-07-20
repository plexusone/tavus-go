package omniavatar

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	tavussdk "github.com/plexusone/tavus-go"
	"github.com/plexusone/tavus-go/api"

	"github.com/plexusone/omniavatar-core/render"
)

// RenderConfig configures the Tavus render (batch video generation) provider.
type RenderConfig struct {
	// APIKey is the Tavus API key.
	// Required.
	APIKey string

	// BaseURL is the Tavus API base URL.
	// Default: https://tavusapi.com
	BaseURL string

	// ReplicaID is the default Tavus replica used when
	// GenerateRequest.AvatarID is empty.
	ReplicaID string

	// HTTPClient is an optional custom HTTP client, used for both API
	// calls and video downloads.
	HTTPClient *http.Client
}

// RenderProvider implements render.Provider for Tavus video generation.
//
// Tavus has no audio upload API, so RenderProvider does not implement
// render.AudioUploader; callers must supply a publicly fetchable
// GenerateRequest.AudioURL (.wav or .mp3).
type RenderProvider struct {
	sdk        *tavussdk.Client
	replicaID  string
	httpClient *http.Client
}

// Compile-time interface check.
var _ render.Provider = (*RenderProvider)(nil)

// NewRenderProvider creates a Tavus render provider.
func NewRenderProvider(cfg RenderConfig) (*RenderProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("%w: APIKey is required", render.ErrInvalidConfig)
	}

	opts := []tavussdk.Option{tavussdk.WithAPIKey(cfg.APIKey)}
	if cfg.BaseURL != "" {
		opts = append(opts, tavussdk.WithBaseURL(cfg.BaseURL))
	}
	if cfg.HTTPClient != nil {
		opts = append(opts, tavussdk.WithHTTPClient(cfg.HTTPClient))
	} else {
		opts = append(opts, tavussdk.WithTimeout(30*time.Second))
	}

	sdk, err := tavussdk.NewClient(opts...)
	if err != nil {
		return nil, render.NewProviderError("tavus", "new_render_provider", err)
	}

	return &RenderProvider{
		sdk:        sdk,
		replicaID:  cfg.ReplicaID,
		httpClient: cfg.HTTPClient,
	}, nil
}

// Name returns the provider name.
func (p *RenderProvider) Name() string { return "tavus" }

// Generate submits a video generation job to Tavus.
//
// GenerateRequest.AvatarID maps to the Tavus replica ID. Width and Height
// are not supported by the Tavus API and are ignored. Background maps
// best-effort: Type "video" sets background_source_url; Type "image" and
// "color" are unsupported. Extensions: "fast" (bool), "callback_url".
func (p *RenderProvider) Generate(ctx context.Context, req render.GenerateRequest) (*render.Job, error) {
	if req.AvatarID == "" {
		req.AvatarID = p.replicaID
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	apiReq := &api.CreateVideoRequest{
		ReplicaID: req.AvatarID,
	}
	if req.AudioURL != "" {
		u, err := url.Parse(req.AudioURL)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid AudioURL: %w", render.ErrInvalidRequest, err)
		}
		apiReq.AudioURL = api.NewOptURI(*u)
	} else {
		apiReq.Script = api.NewOptString(req.Script)
	}
	if req.Title != "" {
		apiReq.VideoName = api.NewOptString(req.Title)
	}
	if req.Background != nil && req.Background.Type == "video" {
		u, err := url.Parse(req.Background.Value)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid Background.Value: %w", render.ErrInvalidRequest, err)
		}
		apiReq.BackgroundSourceURL = api.NewOptURI(*u)
	}
	if req.GetBool("fast", false) {
		apiReq.Fast = api.NewOptBool(true)
	}
	if callbackURL := req.GetString("callback_url", ""); callbackURL != "" {
		u, err := url.Parse(callbackURL)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid callback_url: %w", render.ErrInvalidRequest, err)
		}
		apiReq.CallbackURL = api.NewOptURI(*u)
	}

	resp, err := p.sdk.Videos().Create(ctx, apiReq)
	if err != nil {
		return nil, render.NewProviderError("tavus", "generate", err)
	}
	if !resp.VideoID.Set {
		return nil, render.NewProviderError("tavus", "generate",
			fmt.Errorf("%w: create response missing video_id", render.ErrProviderUnavailable))
	}

	return &render.Job{ID: resp.VideoID.Value, Provider: "tavus"}, nil
}

// Status returns the current status of a generation job.
func (p *RenderProvider) Status(ctx context.Context, jobID string) (*render.JobStatus, error) {
	video, err := p.sdk.Videos().Get(ctx, jobID)
	if err != nil {
		return nil, render.NewProviderError("tavus", "status", err)
	}
	return videoToStatus(jobID, video), nil
}

// Download streams the completed video to dst. It prefers the download
// URL and falls back to the stream URL.
func (p *RenderProvider) Download(ctx context.Context, jobID string, dst io.Writer) error {
	status, err := p.Status(ctx, jobID)
	if err != nil {
		return err
	}
	if status.State != render.JobStateCompleted || status.VideoURL == "" {
		return fmt.Errorf("%w: job %s is %s", render.ErrJobNotCompleted, jobID, status.State)
	}

	if err := render.DownloadURL(ctx, p.httpClient, status.VideoURL, dst); err != nil {
		return render.NewProviderError("tavus", "download", err)
	}
	return nil
}

// videoToStatus converts a Tavus Video to a normalized JobStatus.
func videoToStatus(jobID string, video *api.Video) *render.JobStatus {
	status := &render.JobStatus{ID: jobID}
	if video.VideoID.Set {
		status.ID = video.VideoID.Value
	}
	if video.Status.Set {
		status.State = mapVideoState(video.Status.Value)
		status.RawStatus = string(video.Status.Value)
	} else {
		status.State = render.JobStateProcessing
	}
	switch {
	case video.DownloadURL.Set:
		status.VideoURL = video.DownloadURL.Value.String()
	case video.StreamURL.Set:
		status.VideoURL = video.StreamURL.Value.String()
	}
	return status
}

// mapVideoState maps Tavus video statuses to normalized states.
func mapVideoState(s api.VideoStatus) render.JobState {
	switch s {
	case api.VideoStatusQueued:
		return render.JobStatePending
	case api.VideoStatusGenerating:
		return render.JobStateProcessing
	case api.VideoStatusReady:
		return render.JobStateCompleted
	case api.VideoStatusError, api.VideoStatusDeleted:
		return render.JobStateFailed
	default:
		// Unknown states stay non-terminal so pollers keep waiting.
		return render.JobStateProcessing
	}
}
