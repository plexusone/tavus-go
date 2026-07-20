// Package omniavatar provides an OmniAvatar render provider backed by the
// Tavus video generation API.
//
// It implements render.Provider (Generate/Status/Download) for
// audio-driven talking-head generation using Tavus replicas. Tavus has no
// audio upload API, so it does not implement render.AudioUploader — supply
// a publicly fetchable GenerateRequest.AudioURL.
//
// The adapter is constructor-based and depends only on omniavatar-core.
// The batteries omniavatar package registers it; import that to use it by
// name:
//
//	renderer, err := omniavatar.GetRenderProvider("tavus",
//	    omniavatar.WithAPIKey(os.Getenv("TAVUS_API_KEY")),
//	    omniavatar.WithExtension("replica_id", replicaID))
//
// Or construct directly:
//
//	p, err := tavusomni.NewRenderProvider(tavusomni.RenderConfig{
//	    APIKey: os.Getenv("TAVUS_API_KEY"),
//	})
//
// The real-time (live) Tavus provider lives in the batteries omniavatar
// package (github.com/plexusone/omniavatar/providers/tavus), because its
// LiveKit integration depends on that package.
package omniavatar
