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
	go func() { err = p.UpstreamDiscover(srv.Context(), src, c) }()
	for info := range c {
		if err := srv.Send(info); err != nil {
			glog.V(2).Infof("error sending info downstream: %v", err)
		}
	}
	return err
}

const chunkWriterMaxBufSize = 12288

type chunkWriter struct {
	Server scache.ChunkStore_GetChunkServer
	Blob   *scache.Blob
	Range  *scache.ByteRange

	moreThanOne bool
}

func (w *chunkWriter) Write(b []byte) (n int, err error) {
	bufToSend := func() []byte {
		switch {
		case len(b) > chunkWriterMaxBufSize:
			return b[:chunkWriterMaxBufSize]
		default:
			return b
		}
	}()

	c := scache.Chunk{Data: bufToSend}

	if w.moreThanOne {
		c.Blob = w.Blob
		c.Range = w.Range
	}
	w.moreThanOne = true

	return len(bufToSend), w.Server.Send(&c)
}

func (r *Registry) GetChunk(
	req *scache.ChunkRequest,
	srv scache.ChunkStore_GetChunkServer,
) error {

	if req == nil || req.Blob == nil || req.Blob.SourceType == "" {
		return BadChunkRequest
	}

	p := r.typeMap[req.Blob.SourceType]
	if p == nil {
		return TypeNotHandled
	}

	return p.UpstreamGetChunk(srv.Context(), req, &chunkWriter{
		Server: srv,
		Blob:   req.Blob,
		Range:  req.Range,
	})
}
