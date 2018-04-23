package registry

import (
	"context"

	"github.com/gogo/protobuf/proto"
)

// Upstream that manages a set of source types. It is only required to act as
// an Upstream for source types that are registered with it.
type Provider interface {
	Upstream

	// Returns the protocol and then the uri that identifies this Provider's
	// underlying connection. Providers with the same Addr() should be
	// essentially identical other than handled types.
	Addr() (string, string)

	// Returns the options that this Provider was configured with, if any.
	Opts() proto.Message

	// Registers and deregisters source types that this Provider is responsible
	// for handling. Successful invocations should affect the result of Num().
	StartHandlingType(ctx context.Context, sourceType string) error
	StopHandlingType(ctx context.Context, sourceType string) error

	// Returns the number of handled types. When this reaches 0, any underlying
	// connection can be closed, although this would probably affect in progress
	// Upstream calls, but this is up to the implementation.
	Num() int
}
