package registry

import (
	"context"
	"io"

	"github.com/jonnrb/scache/proto/scache"
)

// Upstream that manages a set of source types. It is only required to act as
// an Upstream for source types that are registered with it.
type Provider interface {
	Upstream

	// Returns the protocol and then the uri that identifies this Provider's
	// underlying connection. Providers with the same Addr() should be
	// essentially identical other than handled types.
	Addr() (string, string)

	// Registers and deregisters source types that this Provider is responsible
	// for handling. Successful invocations should affect the result of Num().
	StartHandlingType(ctx context.Context, sourceType string) error
	StopHandlingType(ctx context.Context, sourceType string) error

	// Returns the number of handled types. When this reaches 0, any underlying
	// connection can be closed, although this would probably affect in progress
	// Upstream calls, but this is up to the implementation.
	Num() int
}

// Generic interface used by the registry to talk to upstream providers.
type Upstream interface {
	// Discovers blobs in |src| from upstream, outputting on |out| until the
	// provided context expires or there is an error, which is relayed unless
	// it occurs within normal operation (e.g. upstream closes stream
	// gracefully). |out| will be closed in any return scenario.
	UpstreamDiscover(
		ctx context.Context,
		src *scache.Source,
		out chan<- *scache.DiscoveryInfo,
	) error

	// During discovery, this writes the chunk specified by |req| into |w|.
	UpstreamGetChunk(
		ctx context.Context,
		req *scache.ChunkRequest,
		w io.Writer,
	) error
}
