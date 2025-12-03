# genproto

Generates Go and JavaScript code from `.proto` files.

## Usage

From the repository root:

```bash
go run ./cmd/tools/genproto
```

Options:

```bash
go run ./cmd/tools/genproto -v          # verbose output
go run ./cmd/tools/genproto -dry-run    # preview without executing
```

For JetBrains IDEs (GoLand, IntelliJ IDEA), use the "gen-proto" run configuration.

## Output

- **Go code**: `pkg/gen/*.pb.go` - Generated from proto definitions
- **JavaScript code**: `ui/ui-gen/proto/*_pb.js` - For grpc-web clients

## Prerequisites

The following tools must be installed and available on your PATH:

```bash
# Go protobuf tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Smart Core tools
go install github.com/smart-core-os/protomod@latest
go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-router@latest
go install github.com/smart-core-os/sc-golang/cmd/protoc-gen-wrapper@latest
```

Also required:
- `protoc` - Protocol Buffers compiler (https://grpc.io/docs/protoc-installation/)
- `protoc-gen-grpc-web` - JavaScript/TypeScript gRPC-Web plugin
- `yarn` - JavaScript package manager

