/*
Package source provides a wrapper around `registry.Upstream` that handles the
protocol boilerplate and exposes a connection-oriented API. A source is
connected to via a `registry.Upstream` and blobs can be opened as `io.Reader`s
while the connection is open. When the source is no longer needed, it can be
closed and associated resources reclaimed.
*/
package source
