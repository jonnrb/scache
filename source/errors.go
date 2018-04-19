package source

import "errors"

var (
	InflationMessageNotReceivedFirst = errors.New(
		"first DiscoveryInfo did not contain inflated source")
	BlobNotFound    = errors.New("blob not found")
	SeekOutOfBounds = errors.New("would seek out of bounds of the blob")
)
