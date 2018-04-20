package main

import (
	"flag"
	"net"

	"github.com/golang/glog"
	"github.com/jonnrb/scache/impl/passthrough"
	"github.com/jonnrb/scache/proto/scache"
	"github.com/jonnrb/scache/registry"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	srv := grpc.NewServer()

	reg := registry.New()
	scache.RegisterProviderRegistryServer(srv, reg)
	scache.RegisterChunkStoreServer(srv, reg)
	scache.RegisterProviderServer(srv, reg)

	cache := passthrough.Service{Registry: reg}
	scache.RegisterCacheServer(srv, &cache)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		glog.Fatal(err)
	}

	if err := srv.Serve(lis); err != nil {
		glog.Fatal(err)
	}
}
