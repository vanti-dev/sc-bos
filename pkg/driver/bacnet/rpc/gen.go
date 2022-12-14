package rpc

// PREREQUISITE: protoc-gen-router is on PATH, i.e. `go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-router@latest`
// PREREQUISITE: protoc-gen-wrapper is on PATH, i.e. `go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-wrapper@latest`
// You will need to copy the files from {root}/pkg/trait/bacnet/ into this package after this is run
//go:generate protoc -I=../../../.. --go_out=paths=source_relative:../../../.. --go-grpc_out=paths=source_relative:../../../.. --wrapper_out=paths=source_relative:../../../.. --router_out=paths=source_relative:../../../.. pkg/driver/bacnet/rpc/bacnet.proto
