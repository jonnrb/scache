package provider

import (
	"context"
	"io"
	"sync"

	"github.com/jonnrb/scache/proto/scache"

	"google.golang.org/grpc"
)

type GRPCProvider struct {
	uri  string
	conn *grpc.ClientConn
	cli  scache.ProviderClient

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
//
func (p *GRPCProvider) establishConn() error {
	if p.conn != nil {
		return nil
	}

	if conn, err := grpc.Dial(p.uri); err != nil {
		return err
	} else {
		p.conn = conn
	}

	p.cli = scache.NewProviderClient(p.conn)
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
	if _, err := p.cli.SupportsType(ctx, &s); err != nil {
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

func (p *GRPCProvider) Discover(
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

		return p.cli.Discover(ctx, src)
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
