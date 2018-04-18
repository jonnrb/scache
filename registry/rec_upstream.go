package registry

import (
	"context"
	"io"

	"github.com/jonnrb/scache/proto/scache"
)

func (r *Registry) UpstreamDiscover(
	ctx context.Context,
	src *scache.Source,
	out chan<- *scache.DiscoveryInfo,
) error {

	if src.SourceType == "" {
		return NoTypeProvided
	}

	p := func() Provider {
		r.mu.RLock()
		defer r.mu.RUnlock()
		return r.typeMap[src.SourceType]
	}()
	if p == nil {
		return TypeNotHandled
	}

	return p.UpstreamDiscover(ctx, src, out)
}

func (r *Registry) UpstreamGetChunk(
	ctx context.Context,
	req *scache.ChunkRequest,
	w io.Writer,
) error {

	if req == nil || req.Blob == nil || req.Blob.SourceType == "" {
		return BadChunkRequest
	}

	p := r.typeMap[req.Blob.SourceType]
	if p == nil {
		return TypeNotHandled
	}

	return p.UpstreamGetChunk(ctx, req, w)
}
