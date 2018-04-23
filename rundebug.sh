#!/bin/bash

V=2

function kill_servers() {
  [ -n "$scache_pid" ] && kill $scache_pid
  [ -n "$scache_gw_pid" ] && kill $scache_gw_pid
  [ -n "$provider_pid" ] && kill $provider_pid
}
trap kill_servers INT

go run ./cmd/scache-passthrough/*.go -logtostderr -v=$V &
scache_pid=$!

pushd ../scache-torrent-provider >/dev/null
node index.js &
provider_pid=$!
popd >/dev/null

sleep 1

go run ./cmd/scache-http-gateway/*.go -logtostderr -v=$V &
scache_gw_pid=$!

go run ./cmd/scache-register-backend/*.go -sourceType torrent

archlinux='{"source":{"sourceType":"torrent","link":"306583088ef59be1698dec4c8967cafb3894e256"}}'
sintel='{"source":{"sourceType":"torrent","link":"08ada5a7a6183aae1e09d831df6748d566095a10"}}'
if [ -n "$1" ]; then
  usertorrent="{\"source\":{\"sourceType\":\"torrent\",\"link\":\"$1\"}}"
  grpcurl -plaintext -d $usertorrent localhost:8080 scache.Cache/AddSource
else
  grpcurl -plaintext -d $sintel localhost:8080 scache.Cache/AddSource
fi

wait $scache_pid
wait $scache_gw_pid
wait $provider_pid
