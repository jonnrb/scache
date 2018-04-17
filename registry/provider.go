package registry

import (
	"context"

	"github.com/jonnrb/scache/proto/scache"
)

// Generic interface used by the registry to talk to upstream providers.
//
type Provider interface {
	Addr() (string, string)

	StartHandlingType(ctx context.Context, sourceType string) error
	StopHandlingType(ctx context.Context, sourceType string) error
	Num() int

	// Discovers blobs in |src| from upstream, outputting on |out| until the
	// provided context expires or there is an error, which is relayed unless
	// it occurs within normal operation (e.g. upstream closes stream
	// gracefully). |out| will be closed in any return scenario.
	//
	Discover(ctx context.Context, src *scache.Source, out chan<- *scache.DiscoveryInfo) error
}
