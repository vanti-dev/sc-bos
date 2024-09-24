package gen

//go:generate protoc -I ../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. enrollment.proto lighting_test.proto

// PREREQUISITE: protomod is on PATH, i.e. `go install github.com/smart-core-os/protomod`
// PREREQUISITE: protoc-gen-router is on PATH, i.e. `go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-router@latest`
// PREREQUISITE: protoc-gen-wrapper is on PATH, i.e. `go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-wrapper@latest`
// You will need to copy the files from {root}/pkg/trait/gen/ into this package after this is run
//go:generate protomod protoc -- -I ../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --router_out=../.. --wrapper_out=../.. access.proto alerts.proto axiomxa.proto button.proto color.proto dali.proto history.proto meter.proto mqtt.proto services.proto status.proto temperature.proto udmi.proto waste.proto
//go:generate protomod protoc -- -I ../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --wrapper_out=../.. devices.proto hub.proto lighting_test.proto tenants.proto
