# API (Resource Server)
The resource server is implemented in Go.

## HTTP Server
It exposes a REST endpoint at `/api/v1/test`,
supporting the `GET` and `PUT` methods with string data. It requires authentication from
the client in the form of an `Authorization: Bearer <token>` header. It verifies this token
using Microsoft's public keys. It will also check for the appropriate *scopes* and *roles*.
The *scopes* are requested when the token is obtained, so the client chooses what scopes
it wants. Scopes are only available for user tokens.
The *roles* are a property of the user/app, and may for example be derived from what
AD Groups the user is a member of. They are available for both user tokens and app tokens.

The API server requires that the token has the `Test.User` and/or `Test.Admin` roles.
To `GET` data, the `Test.Read` scope is required. To `PUT` data, the `Test.Write` scope is
required.

## gRPC server
The server also hosts a plaintext gRPC server. The service definition can be found under `proto/test.proto`.
The authorization requirements are the same as the HTTP server:
it expects an `Authorization: Bearer <token>` in the call metadata.

## Running
Simply run the go program at the root of the repository. No secrets are required.

The files in `static`, including  the web app, are also hosted here for convenience,
but hosting them on the same server is not required as long as the API server's CORS
config is correct.
