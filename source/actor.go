package source

import (
	"context"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
	"github.com/jonnrb/scache/registry"
	"golang.org/x/sync/errgroup"
)

type BlobObserver interface {
	BlobAdded(b *scache.Blob)
	BlobRemoved(b *scache.Blob)
}

type DiscoveryHandler interface {
	ReceiveInflatedSource(src *scache.Source) error
	BlobObserver
}

type CallbackDiscoveryHandler struct {
	ReceiveInflatedSourceCb func(*scache.Source) error
	BlobAddedCb             func(*scache.Blob)
	BlobRemovedCb           func(*scache.Blob)
}

func (h *CallbackDiscoveryHandler) ReceiveInflatedSource(
	src *scache.Source,
) error {

	switch cb := h.ReceiveInflatedSourceCb; cb {
	case nil:
		return nil
	default:
		return cb(src)
	}
}

func (h *CallbackDiscoveryHandler) BlobAdded(b *scache.Blob) {
	if cb := h.BlobAddedCb; cb != nil {
		cb(b)
	}
}

func (h *CallbackDiscoveryHandler) BlobRemoved(b *scache.Blob) {
	if cb := h.BlobRemovedCb; cb != nil {
		cb(b)
	}
}

func ConnectToUpstream(
	ctx context.Context,
	src *scache.Source,
	u registry.Upstream,
	dh DiscoveryHandler,
) error {
	g, ctx := errgroup.WithContext(ctx)

	diC := make(chan *scache.DiscoveryInfo)
	g.Go(func() error {
		// First thing sent down the pipe should be the inflated scache.Source.
		var inflated *scache.Source
		select {
		case initial := <-diC:
			switch i := initial.Info.(type) {
			case *scache.DiscoveryInfo_Inflated:
				inflated = i.Inflated
			default:
				return InflationMessageNotReceivedFirst
			}
		case <-ctx.Done():
			return ctx.Err()
		}

		// Allow handler to short-circuit everything if it didn't like the
		// inflated scache.Source.
		if err := dh.ReceiveInflatedSource(inflated); err != nil {
			return err
		}

		for di := range diC {
			switch b := di.Info.(type) {
			case *scache.DiscoveryInfo_BlobAdded:
				dh.BlobAdded(b.BlobAdded)
			case *scache.DiscoveryInfo_BlobRemoved:
				dh.BlobRemoved(b.BlobRemoved)
			default:
				glog.V(2).Infof("got unexpected DiscoveryInfo: %+v", *di)
			}
		}

		return nil
	})

	g.Go(func() error { return u.UpstreamDiscover(ctx, src, diC) })

	return g.Wait()
}
