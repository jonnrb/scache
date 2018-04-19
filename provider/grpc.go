package provider

import (
	"context"
	"io"
	"sync"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
	"google.golang.org/grpc"
)

type GRPCProvider struct {
	uri string

	conn *grpc.ClientConn
	pCli scache.ProviderClient
	cCli scache.ChunkStoreClient

	mu    sync.RWMutex
	types map[string]bool
}

func NewGRPCProvider(uri string) *GRPCProvider {
	return &GRPCProvider{
		uri:   uri,
		types: make(map[string]bool),
	}
}

func (p *GRPCProvider) Addr() (string, string) {
	return "grpc", p.uri
}

// Write mutex should be held.
func (p *GRPCProvider) establishConn() error {
	if p.conn != nil {
		return nil
	}

	if conn, err := grpc.Dial(p.uri); err != nil {
		return err
	} else {
		p.conn = conn
	}

	p.pCli = scache.NewProviderClient(p.conn)
	p.cCli = scache.NewChunkStoreClient(p.conn)
	return nil
}

func (p *GRPCProvider) StartHandlingType(
	ctx context.Context,
	sourceType string,
) error {

	if err := func() error {
		p.mu.Lock()
		defer p.mu.Unlock()

		if alreadyHandlesType := p.types[sourceType]; alreadyHandlesType {
			return nil
		}

		p.types[sourceType] = true

		if err := p.establishConn(); err != nil {
			return err
		}

		return nil

	}(); err != nil {
		return err
	}

	s := scache.Source{SourceType: sourceType}
	if _, err := p.pCli.SupportsType(ctx, &s); err != nil {
		return err
	}

	return nil
}

func (p *GRPCProvider) StopHandlingType(
	ctx context.Context,
	sourceType string,
) error {

	p.mu.Lock()
	defer p.mu.Unlock()

	if alreadyHandlesType := p.types[sourceType]; !alreadyHandlesType {
		return nil
	}

	delete(p.types, sourceType)

	if len(p.types) == 0 {
		if err := p.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (p *GRPCProvider) Num() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.types)
}

func (p *GRPCProvider) UpstreamDiscover(
	ctx context.Context,
	src *scache.Source,
	out chan<- *scache.DiscoveryInfo,
) error {

	defer close(out)

	cli, err := func() (scache.Provider_DiscoverClient, error) {
		p.mu.RLock()
		defer p.mu.RUnlock()

		if handlesType := p.types[src.SourceType]; !handlesType {
			return nil, TypeNotAvailable
		}

		return p.pCli.Discover(ctx, src)
	}()
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

func (p *GRPCProvider) UpstreamGetChunk(
	ctx context.Context,
	req *scache.ChunkRequest,
	w io.Writer,
) error {

	if req == nil || req.Blob == nil {
		return ChunkRequestMissingBlob
	}

	cli, err := func() (scache.ChunkStore_GetChunkClient, error) {
		p.mu.RLock()
		defer p.mu.RUnlock()

		if handlesType := p.types[req.Blob.SourceType]; !handlesType {
			return nil, TypeNotAvailable
		}

		return p.cCli.GetChunk(ctx, req)
	}()
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
