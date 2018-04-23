package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/glog"
	"github.com/jonnrb/scache/proto/scache"
	"github.com/jonnrb/scache/registry"
	"google.golang.org/grpc"
)

var pm = jsonpb.Marshaler{
	Indent: "  ",
}

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	sc := scache.NewCacheClient(cc)
	u := registry.NewGRPCUpstream(cc)

	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) < 1 {
			panic("wtf")
		}
		trimPath := string(([]rune(r.URL.Path))[1:])
		switch parts := strings.Split(trimPath, "/"); len(parts) {
		case 1:
			glog.V(2).Info("got request for /")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			res, err := sc.ListSources(ctx, &scache.ListSourcesRequest{})
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			listSourceResponse(w, res)

		case 2:
			sourceType, link := parts[0], parts[1]
			if sourceType == "" || link == "" {
				http.Error(w, "not found", 404)
				return
			}

			glog.V(2).Infof("got request for /%v/%v", sourceType, link)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			res, err := sc.ListBlobs(ctx, &scache.ListBlobsRequest{
				Filter: &scache.FilterExpression{
					Term: []*scache.FilterExpression_Term{
						&scache.FilterExpression_Term{
							SourceType: sourceType,
							Link:       link,
						},
					},
				},
			})
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			listBlobsResponse(w, sourceType, link, res)

		case 3:
			sourceType, link, blobId := parts[0], parts[1], parts[2]
			if sourceType == "" || link == "" || blobId == "" {
				http.Error(w, "not found", 404)
				return
			}

			glog.V(2).Infof("got request for /%v/%v/%v", sourceType, link, blobId)

			serveBlob(u, r, w, sourceType, link, blobId)

		default:
			http.Error(w, "not found", 404)
		}
	}))
}
