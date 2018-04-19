package source

import (
	"context"

	"github.com/jonnrb/scache/proto/scache"
	"github.com/jonnrb/scache/registry"
)

type Config struct {
	// Optional. Set to receive notifications.
	Observer BlobObserver
}

type Conn struct {
	state state

	u      registry.Upstream
	ctx    context.Context
	cancel context.CancelFunc
	errC   chan error
}

// The ctx passed in only scopes the initial source discovery, after which,
// this function will return. To close the Conn, Close() must be called on it.
func Open(
	ctx context.Context,
	src *scache.Source,
	u registry.Upstream,
	cfg *Config,
) (*Conn, error) {

	connCtx, cancel := context.WithCancel(context.Background())

	s := state{
		connected: make(chan struct{}), // Closed upon connection.
		blobs:     make(map[string]*scache.Blob),
	}

	go func() {
		select {
		case <-s.connected:
			return
		case ctx.Done():
			cancel()
		}
	}()

	errC := make(chan error, 1)
	go func() {
		errC <- ConnectToUpstream(connCtx, src, u, &s)
		cancel()
	}()

	select {
	case err := <-errC:
		// Upstream context err supplants ConnectToUpstream err.
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		return nil, err

	case <-s.connected:
		go func() {
			switch err := <-errC; err {
			case context.Canceled, nil:
			default:
				s.SetError(err)
			}
		}()

		return &Conn{
			u:      u,
			ctx:    connCtx,
			cancel: cancel,
		}, nil
	}
}

func (c *Conn) Close() error {
	if err := c.state.GetError(); err != nil {
		return err
	}
	c.cancel()
	return nil
}

// Do not edit the returned object.
func (c *Conn) InflatedSource() *scache.Source {
	return c.state.inflated
}

// Do not edit any of the returned objects.
func (c *Conn) Blobs() []*scache.Blob {
	c.state.mu.RLock()
	defer c.state.mu.RUnlock()

	var l []*scache.Blob
	for _, b := range c.state.blobs {
		l = append(l, b)
	}

	return l
}
