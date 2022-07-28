package gen

//go:generate protoc -I ../../proto test.proto enrollment.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:.
