package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
	"github.com/jonnrb/scache/registry"
	"github.com/jonnrb/scache/source"
)

func listSourceResponse(w io.Writer, res *scache.SourceList) {
	listWriter := func(w io.Writer) error {
		for _, src := range res.Source {
			title := func() string {
				if src.Metadata != nil && len(src.Metadata.Name) != 0 {
					return src.Metadata.Name
				} else {
					return fmt.Sprintf("%s: %s", src.SourceType, src.Link)
				}
			}()
			link := fmt.Sprintf(
				"/%s/%s",
				url.PathEscape(src.SourceType),
				url.PathEscape(src.Link),
			)
			if err := writeListLink(w, title, link); err != nil {
				return err
			}
		}
		return nil
	}

	pt, err := pm.MarshalToString(res)
	if err != nil {
		panic(err)
	}

	err = writeChain(w,
		headerWriter("All Sources"),

		startList,
		listWriter,
		endList,

		startCode,
		stringWriter(pt),
		endCode,

		writeFooter,
	)
	if err != nil {
		panic(err)
	}
}

func listBlobsResponse(
	w io.Writer,
	sourceType, link string,
	res *scache.BlobList,
) {

	listWriter := func(w io.Writer) error {
		for _, blob := range res.Blob {
			title := func() string {
				if blob.Metadata != nil && len(blob.Metadata.Name) != 0 {
					return blob.Metadata.Name
				} else {
					return fmt.Sprintf("%s: %s [%s]",
						blob.SourceType, blob.Link, blob.BlobId)
				}
			}()
			link := fmt.Sprintf(
				"/%s/%s/%s",
				url.PathEscape(blob.SourceType),
				url.PathEscape(blob.Link),
				url.PathEscape(blob.BlobId),
			)
			if err := writeListLink(w, title, link); err != nil {
				return err
			}
		}
		return nil
	}

	pt, err := pm.MarshalToString(res)
	if err != nil {
		panic(err)
	}

	err = writeChain(w,
		headerWriter(fmt.Sprintf("%s: %s", sourceType, link)),

		startList,
		listWriter,
		endList,

		startCode,
		stringWriter(pt),
		endCode,

		writeFooter,
	)
	if err != nil {
		panic(err)
	}
}

func serveBlob(
	u registry.Upstream,
	r *http.Request,
	w http.ResponseWriter,
	sourceType, link, blobId string,
) {
	src := scache.Source{
		SourceType: sourceType,
		Link:       link,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	go func() {
		select {
		case <-r.Cancel:
			cancel()
		case <-ctx.Done():
		}
	}()

	mu := make(chan struct{})
	discovery := make(chan *scache.Blob)
	waiter := source.CallbackDiscoveryHandler{
		BlobAddedCb: func(b *scache.Blob) {
			if b.BlobId == blobId {
				select {
				case _, ok := <-mu:
					if ok {
						discovery <- b
						close(discovery)
					}
				}
			}
		},
	}

	conn, err := func() (*source.Conn, error) {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		return source.Open(ctx, &src, u, &source.Config{Observer: &waiter})
	}()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer conn.Close()

	err = func() error {
		defer close(mu)

		select {
		case <-conn.Done():
			return conn.Err()
		case mu <- struct{}{}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	blob := <-discovery

	br, err := conn.NewReader(blob)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	name := func() string {
		var name string
		if blob.Metadata != nil {
			name = blob.Metadata.Name
		}
		if name == "" {
			name = blob.BlobId
		}
		return name
	}()

	glog.V(2).Infof("serving content for /%v/%v/%v", sourceType, link, blobId)
	http.ServeContent(w, r, name, time.Time{}, br)
	glog.V(2).Infof("finished serving content for /%v/%v/%v", sourceType, link, blobId)
}
