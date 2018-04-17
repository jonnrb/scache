package registry

import (
	"context"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
)

var (
	supportsTypeRes = scache.SupportsTypeResponse{}
)

func (r *Registry) SupportsType(
	ctx context.Context,
	src *scache.Source,
) (*scache.SupportsTypeResponse, error) {

	if src.SourceType == "" {
		return nil, NoTypeProvided
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	p := r.typeMap[src.SourceType]
	if p == nil {
		return nil, TypeNotHandled
	}

	return &supportsTypeRes, nil
}

func (r *Registry) Discover(
	src *scache.Source,
	srv scache.Provider_DiscoverServer,
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

	var (
		c   = make(chan *scache.DiscoveryInfo)
		err error
	)
	go func() { err = p.Discover(srv.Context(), src, c) }()
	for info := range c {
		if err := srv.Send(info); err != nil {
			glog.V(2).Infof("error sending info downstream: %v", err)
		}
	}
	return err
}
