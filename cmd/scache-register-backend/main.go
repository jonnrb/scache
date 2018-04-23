package main

import (
	"flag"
	"time"

	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/jonnrb/scache/proto/scache"
	"google.golang.org/grpc"
)

var (
	reg   = flag.String("reg", "localhost:8080", "address of the registry")
	regPt = flag.Bool("pt-reg", true, "registry is insecure (no TLS)")
	be    = flag.String("be", "localhost:8081", "address of the  backend")
	bePt  = flag.Bool("pt-be", true, "backend is insecure (no TLS)")
	st    = flag.String("sourceType", "", "source type to register")
)

func main() {
	flag.Parse()

	if *st == "" {
		panic("no sourceType set")
	}

	opts := func() *any.Any {
		if *bePt {
			opts, err := ptypes.MarshalAny(&scache.GRPCProviderOpts{
				DialInsecure: *bePt,
			})
			if err != nil {
				panic(err)
			}
			return opts
		}
		return nil
	}()
	spec := scache.ProviderSpec{
		Addr: &scache.ProviderAddress{
			Uri:     *be,
			Proto:   "grpc",
			Options: opts,
		},
		SourceType: []string{*st},
	}

	var dialOpts []grpc.DialOption
	if *regPt {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*reg, dialOpts...)
	if err != nil {
		panic(err)
	}

	cli := scache.NewProviderRegistryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := cli.AddProvider(ctx, &spec)
	_ = res

	if err != nil {
		panic(err)
	}
	println("success")
}
