package provider

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	TypeNotAvailable = status.Error(
		codes.Unavailable, "source type not available")
)
