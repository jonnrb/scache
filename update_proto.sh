#!/bin/bash
# Generates go code for protos used throughout the API server.

cd proto
protoc --go_out=plugins=grpc:. $(find . |grep '\.proto$')
