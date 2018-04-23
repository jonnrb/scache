package registry

import (
	"context"
	"io"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
	"google.golang.org/grpc"
)

// Generic interface used by the registry to talk to upstream providers.
type Upstream interface {
	// Discovers blobs in `src` from upstream, outputting on `out` until the
	// provided context expires or there is an error, which is relayed unless
	// it occurs within normal operation (e.g. upstream closes stream
	// gracefully). `out` will be closed in any return scenario.
	UpstreamDiscover(
		ctx context.Context,
		src *scache.Source,
		out chan<- *scache.DiscoveryInfo,
	) error

	// During discovery, this writes the chunk specified by `req` into `w`.
	UpstreamGetChunk(
		ctx context.Context,
		req *scache.ChunkRequest,
		w io.Writer,
	) error
}

type GRPCUpstream struct {
	providerCli scache.ProviderClient
	chunkCli    scache.ChunkStoreClient
}

func NewGRPCUpstream(cc *grpc.ClientConn) *GRPCUpstream {
	return &GRPCUpstream{
		providerCli: scache.NewProviderClient(cc),
		chunkCli:    scache.NewChunkStoreClient(cc),
	}
}

func (p *GRPCUpstream) UpstreamDiscover(
	ctx context.Context,
	src *scache.Source,
	out chan<- *scache.DiscoveryInfo,
) error {

	defer close(out)

	cli, err := p.providerCli.Discover(ctx, src)
	if err != nil {
		return err
	}

	for {
		d, err := cli.Recv()

		switch err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}

		out <- d
	}
}

func (p *GRPCUpstream) UpstreamGetChunk(
	ctx context.Context,
	req *scache.ChunkRequest,
	w io.Writer,
) error {

	if req == nil || req.Blob == nil {
		return BadChunkRequest
	}

	cli, err := p.chunkCli.GetChunk(ctx, req)
	if err != nil {
		return err
	}

	for {
		c, err := cli.Recv()

		switch err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}

		var n int
		for n < len(c.Data) {
			if m, err := w.Write(c.Data[n:]); err != nil {
				if err := cli.CloseSend(); err != nil {
					glog.Errorf("error closing stream: %v", err)
				}
				return err
			} else {
				n += m
			}
		}
	}
}
