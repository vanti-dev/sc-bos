package gen

//go:generate protoc -I ../../proto test.proto enrollment.proto nodes.proto tenants.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:.
