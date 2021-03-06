package registry

import (
	"context"
	"sync"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/jonnrb/scache/proto/scache"
	"github.com/jonnrb/scache/provider"
)

type Registry struct {
	mu        sync.RWMutex
	typeMap   map[string]Provider
	providers map[string]map[string]Provider
}

var (
	addProviderRes    = scache.AddProviderResponse{}
	removeProviderRes = scache.RemoveProviderResponse{}
)

func New() *Registry {
	return &Registry{
		typeMap:   make(map[string]Provider),
		providers: make(map[string]map[string]Provider),
	}
}

func (r *Registry) providerMapAdd(proto, uri string, p Provider) {
	if r.providers[proto] == nil {
		r.providers[proto] = make(map[string]Provider)
	}
	r.providers[proto][uri] = p
}

func (r *Registry) providerMapGet(proto, uri string) Provider {
	switch {
	case r.providers[proto] != nil:
		return r.providers[proto][uri]
	default:
		return nil
	}
}

func (r *Registry) providerMapDelete(proto, uri string) {
	if r.providers[proto] != nil {
		delete(r.providers[proto], uri)
	}
	if len(r.providers[proto]) == 0 {
		delete(r.providers, proto)
	}
}

func (r *Registry) AddProvider(
	ctx context.Context,
	s *scache.ProviderSpec,
) (*scache.AddProviderResponse, error) {

	// TODO: Shorten this critical section.
	r.mu.Lock()
	defer r.mu.Unlock()

	if s.Addr == nil {
		return nil, NoAddrProvided
	}
	proto, uri, opts := s.Addr.Proto, s.Addr.Uri, s.Addr.Options

	// The default proto is gRPC.
	if proto == "" {
		proto = "grpc"
	}

	if len(s.SourceType) == 0 {
		return nil, NoTypesProvided
	}

	var p Provider
	if p = r.providerMapGet(proto, uri); p == nil {
		switch proto {
		case "grpc":
			gOpts, err := func() (*scache.GRPCProviderOpts, error) {
				switch {
				case opts != nil:
					gOpts := &scache.GRPCProviderOpts{}
					if err := ptypes.UnmarshalAny(opts, gOpts); err != nil {
						return nil, err
					}
					return gOpts, nil
				default:
					return nil, nil
				}
			}()
			if err != nil {
				return nil, err
			}
			p = provider.NewGRPCProvider(uri, gOpts)

		default:
			glog.V(2).Infof("got unknown proto: %q", proto)
			return nil, UnimplementedProto
		}
	}

	// Register this provider for these types, and rollback on failure.
	typesSucc := 0
	rollBack := func(err error) error {
		for _, t := range s.SourceType[:typesSucc] {
			if err := p.StopHandlingType(ctx, t); err != nil {
				return err
			}
		}
		return err
	}
	for _, t := range s.SourceType {
		if err := p.StartHandlingType(ctx, t); err != nil {
			return nil, rollBack(err)
		}
		typesSucc += 1
	}

	// Deregister existing providers of these source types. Don't rollback on
	// failure (log it).
	for _, t := range s.SourceType {
		if other := r.typeMap[t]; other != nil {
			proto, uri := other.Addr()
			if err := other.StopHandlingType(ctx, t); err != nil {
				glog.Errorf(
					"%q (via %v) would not stop handling source type %q: %v",
					uri, proto, t, err,
				)
			}
			if other.Num() == 0 {
				r.providerMapDelete(proto, uri)
			}
		}
		r.typeMap[t] = p
	}
	r.providerMapAdd(proto, uri, p)

	return &addProviderRes, nil
}

func (r *Registry) RemoveProvider(
	ctx context.Context,
	s *scache.ProviderSpec,
) (*scache.RemoveProviderResponse, error) {

	// TODO: Shorten this critical section; should be pretty easy.
	r.mu.Lock()
	defer r.mu.Unlock()

	if s.Addr == nil {
		return nil, NoAddrProvided
	}
	proto, uri := s.Addr.Proto, s.Addr.Uri

	// The default proto is gRPC.
	if proto == "" {
		proto = "grpc"
	}

	if len(s.SourceType) == 0 {
		return nil, NoTypesProvided
	}

	p := r.providerMapGet(proto, uri)
	if p == nil {
		return nil, ProviderNotFound
	}

	// Make sure all requested types are actually handled by this provider. It
	// is a failed precondition if this is not the case. Do nothing on error.
	typesSucc := 0
	rollback := func() {
		for _, t := range s.SourceType[:typesSucc] {
			r.typeMap[t] = p
		}
	}
	for _, t := range s.SourceType {
		if r.typeMap[t] != p {
			rollback()
			return nil, TypeNotHandled
		}
		r.typeMap[t] = nil
		typesSucc += 1
	}

	// Notify the backend that it no longer handles these types. Log any issues.
	for _, t := range s.SourceType {
		if err := p.StopHandlingType(ctx, t); err != nil {
			glog.Errorf(
				"%q (via %v) would not stop handling source type %q: %v",
				uri, proto, t, err,
			)
		}
	}

	if p.Num() == 0 {
		r.providerMapDelete(proto, uri)
	}

	return &removeProviderRes, nil
}

func (r *Registry) ListProviders(
	ctx context.Context,
	_ *scache.ListProvidersRequest,
) (*scache.ListProvidersResponse, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Invert the map.
	m := make(map[Provider][]string)
	for sourceType, p := range r.typeMap {
		m[p] = append(m[p], sourceType)
	}

	var res scache.ListProvidersResponse
	for p, sourceTypes := range m {
		proto, uri := p.Addr()

		opts, err := func() (*any.Any, error) {
			switch opts := p.Opts(); opts {
			case nil:
				return nil, nil
			default:
				return ptypes.MarshalAny(opts)
			}
		}()
		if err != nil {
			glog.Errorf(
				"error marshalling options for %q (via %v): %v",
				uri, proto, err,
			)
		}

		res.Provider = append(res.Provider, &scache.ProviderSpec{
			Addr: &scache.ProviderAddress{
				Uri:     uri,
				Proto:   proto,
				Options: opts,
			},
			SourceType: sourceTypes,
		})
	}
	return &res, nil
}
