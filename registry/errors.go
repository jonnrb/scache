package registry

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	NoAddrProvided     = status.Error(codes.InvalidArgument, "no address provided")
	NoTypesProvided    = status.Error(codes.InvalidArgument, "no source types provided")
	NoTypeProvided     = status.Error(codes.InvalidArgument, "no source type provided")
	UnimplementedProto = status.Error(codes.Unimplemented, "proto not implemented")
	ProviderNotFound   = status.Error(codes.NotFound, "provider not found")
	TypeNotHandled     = status.Error(codes.FailedPrecondition, "provider does not handle type")
)
