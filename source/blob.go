package source

import (
	"io"

	"github.com/jonnrb/scache/proto/scache"
)

type BlobReader struct {
	conn *Conn
	blob *scache.Blob
	pos  int64
}

func (c *Conn) NewReader(b *scache.Blob) (*BlobReader, error) {
	if err := c.state.GetError(); err != nil {
		return nil, err
	}

	c.state.mu.RLock()
	defer c.state.mu.RUnlock()

	if _, hasBlob := c.state.blobs[b.BlobId]; !hasBlob {
		return nil, BlobNotFound
	}

	return &BlobReader{c, b, 0}, nil
}

func min(i, j int) int {
	if i < j {
		return i
	} else {
		return j
	}
}

type chunkWriter []byte

func (w *chunkWriter) Write(p []byte) (n int, err error) {
	copied := min(len(p), len(*w))
	copy(*w, p)
	*w = (*w)[:copied]
	return copied, err
}

func (r *BlobReader) Read(p []byte) (n int, err error) {
	if err := r.conn.state.GetError(); err != nil {
		return 0, err
	}

	n = len(p)
	if r.pos+int64(len(p)) > r.blob.Length {
		n = int(r.blob.Length - r.pos)
	}

	req := scache.ChunkRequest{
		Blob: r.blob,
		Range: &scache.ByteRange{
			Low:  r.pos,
			High: r.pos + int64(n),
		},
	}
	err = r.conn.u.UpstreamGetChunk(r.conn.ctx, &req, (*chunkWriter)(&p))
	return
}

func (r *BlobReader) Seek(offset int64, whence int) (int64, error) {
	if err := r.conn.state.GetError(); err != nil {
		return 0, err
	}

	switch whence {
	case io.SeekStart:
		if offset > r.blob.Length {
			return 0, SeekOutOfBounds
		}
		r.pos = offset

	case io.SeekCurrent:
		switch {
		case -offset > r.pos:
			return 0, SeekOutOfBounds
		case offset > r.blob.Length:
			return 0, SeekOutOfBounds
		case offset+r.pos > r.blob.Length:
			return 0, SeekOutOfBounds
		case offset+r.pos < 0:
			return 0, SeekOutOfBounds
		}
		r.pos += offset

	case io.SeekEnd:
		switch {
		case offset > 0:
			return 0, SeekOutOfBounds
		case -offset > r.blob.Length:
			return 0, SeekOutOfBounds
		}
		r.pos = offset + r.blob.Length
	}

	return r.pos, nil
}
