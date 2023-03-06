package gen

//go:generate protoc -I ../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. test.proto enrollment.proto lighting_test.proto

// PREREQUISITE: protomod is on PATH, i.e. `go install github.com/smart-core-os/protomod`
// PREREQUISITE: protoc-gen-router is on PATH, i.e. `go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-router@latest`
// PREREQUISITE: protoc-gen-wrapper is on PATH, i.e. `go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-wrapper@latest`
// You will need to copy the files from {root}/pkg/trait/gen/ into this package after this is run
//go:generate protomod protoc -- -I ../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --router_out=../.. --wrapper_out=../.. alerts.proto button.proto history.proto meter.proto mqtt.proto priority.proto services.proto udmi.proto
//go:generate protomod protoc -- -I ../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --wrapper_out=../.. devices.proto nodes.proto tenants.proto
