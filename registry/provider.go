package registry

import (
	"context"
	"io"

	"github.com/jonnrb/scache/proto/scache"
)

// Upstream that provides access to
type Provider interface {
	Upstream

	Addr() (string, string)

	StartHandlingType(ctx context.Context, sourceType string) error
	StopHandlingType(ctx context.Context, sourceType string) error
	Num() int
}

// Generic interface used by the registry to talk to upstream providers.
//
type Upstream interface {
	// Discovers blobs in |src| from upstream, outputting on |out| until the
	// provided context expires or there is an error, which is relayed unless
	// it occurs within normal operation (e.g. upstream closes stream
	// gracefully). |out| will be closed in any return scenario.
	//
	UpstreamDiscover(
		ctx context.Context,
		src *scache.Source,
		out chan<- *scache.DiscoveryInfo,
	) error

	// During discovery, this writes the chunk specified by |req| into |w|.
	//
	UpstreamGetChunk(
		ctx context.Context,
		req *scache.ChunkRequest,
		w io.Writer,
	) error
}
