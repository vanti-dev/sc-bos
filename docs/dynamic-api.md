# Dynamic APIs and gRPC

There are a number of limitations with how gRPC and protobufs are implemented in Go that make adding to the API surface
area of a server at runtime difficult. Firstly, the `grpc.Server` type does not provide a way to add new services or
APIs once it has been started. Secondly, protobuf types are registered globally in the `protoregistry` package when they
are imported. Finally, the reflection api is closely tied to the semantics of the `grpc.Server` and `protoregistry`
packages.

SC BOS has a requirement to transparently support APIs that it doesn't know about at compile time, or during boot.
The most pressing of these is for the `proxy` system (used by gateways) to be able to proxy and route RPCs without prior
knowledge of the proto descriptors or service definitions used by the downstream nodes.
The second use case is to allow services (drivers, systems, etc.) to be added to the server at runtime, including
support for their published custom APIs.

Adding transparent support for adjusting the API surface area of a server at runtime is what we call "dynamic APIs".

## Transparency

To achieve transparency for dynamic APIs we need to play nicely with the existing features of a gRPC server.
This includes:

1. Have the gRPC server correctly route RPCs to the correct service implementation.
2. Have the reflection API correctly report both the available RPCs and types for the dynamic APIs.
3. Have our auth policies correctly apply to the dynamic APIs.
4. Have grpc-web work correctly with the dynamic APIs.

## Design

To support dynamic APIs we need two things:

1. Some way to _discover_ the services and types the new APIs provide.
2. Some way to _hook into_ the existing gRPC features to expose this dynamically discovered information.

For the discovery part, we can use the reflection API on the downstream nodes to collect all the information we need.
For the hooks, `grpc.Server` provides a way to register an unknown service handler that we can use to correctly route
RPCs and decode messages for our auth policies to inspect.
Supporting reflection of dynamic APIs is also possible via options on the reflection server, adjusting the type and
service resolvers to include the discovered information.

The Go proto package also provides a dynamicpb package that allows us to construct protobuf messages at runtime from
protobuf descriptors.

A [POC project](https://github.com/vanti-dev/dynamic-api) was written to figure out this design, look there for a more
focused version of a dynamic API proxy.
