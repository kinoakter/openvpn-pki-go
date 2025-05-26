#!/usr/bin/env bash

version='0.0.1'
context_path="."
image="viktorvorobei/ovpn:$version"
platform="linux/amd64,linux/arm64"

#docker build -f ./Dockerfile \
#-t registry.gitlab.com/vpn-tube/vpn-tube-server/ovpn:${version} \
#${context_path}

docker buildx build -f Dockerfile \
  --platform "${platform}" \
  -t $image --push .