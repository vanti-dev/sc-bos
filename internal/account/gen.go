package account

//go:generate sqlc generate
//go:generate protoc --go_out=paths=source_relative:. page_token.proto
