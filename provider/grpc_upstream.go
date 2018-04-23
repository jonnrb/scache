package provider

import (
	"context"
	"io"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
	"google.golang.org/grpc"
)

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
		return ChunkRequestMissingBlob
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
