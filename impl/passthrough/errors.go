package passthrough

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	NoSourceProvided    = status.Error(codes.InvalidArgument, "no source provided")
	SourceAlreadyExists = status.Error(codes.AlreadyExists, "source already exists")
	SourceNotFound      = status.Error(codes.NotFound, "source not found")
	InternalError       = status.Error(codes.Internal, "internal error")
)
