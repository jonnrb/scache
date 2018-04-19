package source

import (
	"sync"

	"github.com/jonnrb/scache/proto/scache"
)

type state struct {
	cfg *Config

	connected chan struct{}
	inflated  *scache.Source

	mu       sync.RWMutex
	blobs    map[string]*scache.Blob
	errState error
}

func (s *state) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errState = err
}

func (s *state) GetError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errState
}

func (s *state) ReceiveInflatedSource(src *scache.Source) error {
	s.inflated = src
	close(s.connected)
	return nil
}

func (s *state) BlobAdded(b *scache.Blob) {
	func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.blobs[b.BlobId] = b
	}()
	if s.cfg != nil && s.cfg.Observer != nil {
		s.cfg.Observer.BlobAdded(b)
	}
}

func (s *state) BlobRemoved(b *scache.Blob) {
	func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.blobs, b.BlobId)
	}()
	if s.cfg != nil && s.cfg.Observer != nil {
		s.cfg.Observer.BlobAdded(b)
	}
}
