package passthrough

import (
	"context"
	"sync"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
	"github.com/jonnrb/scache/registry"
	"github.com/jonnrb/scache/source"
)

type Service struct {
	Registry *registry.Registry

	mu      sync.RWMutex
	sources map[string]map[string]*source.Conn
}

func New(r *registry.Registry) *Service {
	return &Service{
		Registry: r,
		sources:  make(map[string]map[string]*source.Conn),
	}
}

func (s *Service) sourceFor(sourceType, link string) *source.Conn {
	if srcMap, ok := s.sources[sourceType]; !ok {
		return nil
	} else if conn, ok := srcMap[link]; !ok {
		return nil
	} else {
		return conn
	}
}

func (s *Service) AddSource(
	ctx context.Context,
	req *scache.AddSourceRequest,
) (*scache.AddSourceResponse, error) {

	if req.Source == nil {
		return nil, NoSourceProvided
	}
	src := req.Source

	s.mu.Lock()
	defer s.mu.Unlock()

	if c := s.sourceFor(src.SourceType, src.Link); c != nil {
		return nil, SourceAlreadyExists
	}

	conn, err := source.Open(ctx, src, s.Registry, nil)
	if err != nil {
		return nil, err
	}

	srcMap := s.sources[src.SourceType]
	if srcMap == nil {
		srcMap = make(map[string]*source.Conn)
		s.sources[src.SourceType] = srcMap
	}
	srcMap[src.Link] = conn

	return &scache.AddSourceResponse{Source: conn.InflatedSource()}, nil
}

func (s *Service) RemoveSource(
	ctx context.Context,
	req *scache.RemoveSourceRequest,
) (*scache.RemoveSourceResponse, error) {

	if req.Source == nil {
		return nil, NoSourceProvided
	}
	src := req.Source

	s.mu.Lock()
	defer s.mu.Unlock()

	conn := s.sourceFor(src.SourceType, src.Link)
	if conn == nil {
		return nil, SourceNotFound
	}

	if err := conn.Close(); err != nil {
		glog.Errorf("error removing source %+v: %v", src, err)
	}

	srcMap := s.sources[src.SourceType]
	delete(srcMap, src.Link)
	if len(srcMap) == 0 {
		delete(s.sources, src.SourceType)
	}

	return &scache.RemoveSourceResponse{}, nil
}

func (s *Service) ListSources(
	ctx context.Context,
	req *scache.ListSourcesRequest,
) (*scache.SourceList, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	var res scache.SourceList
	for _, srcMap := range s.sources {
		for _, conn := range srcMap {
			if src := conn.InflatedSource(); src != nil {
				res.Source = append(res.Source, conn.InflatedSource())
			} else {
				glog.Error("nil inflated source")
			}
		}
	}
	return &res, nil
}

func (*Service) ObserveSources(
	*scache.ListSourcesRequest,
	scache.Cache_ObserveSourcesServer,
) error {

	panic("implement me")
}

func (s *Service) ListBlobs(
	context.Context,
	*scache.ListBlobsRequest,
) (*scache.BlobList, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	var res scache.BlobList
	for _, srcMap := range s.sources {
		for _, conn := range srcMap {
			res.Blob = append(res.Blob, conn.Blobs()...)
		}
	}
	return &res, nil
}

func (*Service) ObserveBlobs(
	*scache.ListBlobsRequest,
	scache.Cache_ObserveBlobsServer,
) error {

	panic("implement me")
}
