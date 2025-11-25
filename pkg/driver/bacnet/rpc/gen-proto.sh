#!/usr/bin/env bash

go generate
if [ -d ../../../../pkg/trait/rpc ]; then
  sed -i.bak '/rpc "github.com\/smart-core-os\/sc-bos\/pkg\/driver\/bacnet\/rpc"/d' ../../../../pkg/trait/rpc/*.pb.go
  sed -i.bak 's/\brpc\.//g' ../../../../pkg/trait/rpc/*.pb.go
  sed -i.bak 's/[[:<:]]rpc\.//g' ../../../../pkg/trait/rpc/*.pb.go
  sed -i.bak 's/package rpcpb/package rpc/g' ../../../../pkg/trait/rpc/*.pb.go
  mv -f ../../../../pkg/trait/rpc/*.pb.go .
fi
if [ -d ../../../../pkg/trait/rpcpb ]; then
  sed -i.bak '/rpc "github.com\/smart-core-os\/sc-bos\/pkg\/driver\/bacnet\/rpc"/d' ../../../../pkg/trait/rpcpb/*.pb.go
  sed -i.bak 's/\brpc\.//g' ../../../../pkg/trait/rpcpb/*.pb.go
  sed -i.bak 's/[[:<:]]rpc\.//g' ../../../../pkg/trait/rpcpb/*.pb.go
  sed -i.bak 's/package rpcpb/package rpc/g' ../../../../pkg/trait/rpcpb/*.pb.go
  mv -f ../../../../pkg/trait/rpcpb/*.pb.go .
fi

rm -rf ../../../../pkg/trait/


