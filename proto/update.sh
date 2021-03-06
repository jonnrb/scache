#!/bin/bash
# Generates go code for a set of directories' protobufs.

function generate() {
  protoc "--descriptor_set_out=$1_descriptors.protobin" \
    --go_out=plugins=grpc:. \
    $(find "$1" |grep -E "\.proto$")
}

cd $(dirname $(realpath "$0"))
generate ./scache
