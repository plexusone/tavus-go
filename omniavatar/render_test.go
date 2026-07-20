package omniavatar

import (
	"net/url"
	"testing"

	"github.com/plexusone/tavus-go/api"

	"github.com/plexusone/omniavatar-core/render"
)

func TestMapVideoState(t *testing.T) {
	tests := []struct {
		in   api.VideoStatus
		want render.JobState
	}{
		{api.VideoStatusQueued, render.JobStatePending},
		{api.VideoStatusGenerating, render.JobStateProcessing},
		{api.VideoStatusReady, render.JobStateCompleted},
		{api.VideoStatusError, render.JobStateFailed},
		{api.VideoStatusDeleted, render.JobStateFailed},
		{api.VideoStatus("unknown"), render.JobStateProcessing},
	}
	for _, tt := range tests {
		if got := mapVideoState(tt.in); got != tt.want {
			t.Errorf("mapVideoState(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestVideoToStatusPrefersDownloadURL(t *testing.T) {
	download := url.URL{Scheme: "https", Host: "x", Path: "/download.mp4"}
	stream := url.URL{Scheme: "https", Host: "x", Path: "/stream.m3u8"}
	video := &api.Video{
		VideoID:     api.NewOptString("vid-1"),
		Status:      api.NewOptVideoStatus(api.VideoStatusReady),
		DownloadURL: api.NewOptURI(download),
		StreamURL:   api.NewOptURI(stream),
	}

	status := videoToStatus("vid-1", video)
	if status.State != render.JobStateCompleted {
		t.Errorf("State = %q, want %q", status.State, render.JobStateCompleted)
	}
	if status.VideoURL != download.String() {
		t.Errorf("VideoURL = %q, want download URL %q", status.VideoURL, download.String())
	}
}

func TestVideoToStatusFallsBackToStreamURL(t *testing.T) {
	stream := url.URL{Scheme: "https", Host: "x", Path: "/stream.m3u8"}
	video := &api.Video{
		Status:    api.NewOptVideoStatus(api.VideoStatusReady),
		StreamURL: api.NewOptURI(stream),
	}

	status := videoToStatus("vid-1", video)
	if status.ID != "vid-1" {
		t.Errorf("ID = %q, want %q (from jobID when VideoID unset)", status.ID, "vid-1")
	}
	if status.VideoURL != stream.String() {
		t.Errorf("VideoURL = %q, want stream URL %q", status.VideoURL, stream.String())
	}
}

func TestVideoToStatusMissingStatus(t *testing.T) {
	status := videoToStatus("vid-1", &api.Video{})
	if status.State != render.JobStateProcessing {
		t.Errorf("State = %q, want %q (unset status is non-terminal)", status.State, render.JobStateProcessing)
	}
}
