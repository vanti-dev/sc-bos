App-only client example
=======================

This program is an example of an app-only client - an app that accesses a resource server on behalf of itself, rather
than a particular user. It uses the OAuth 2 Client Credentials flow to obtain an access token.
This is a confidential client application, which has a secret, which is how it authenticates itself to Azure AD.

The secret should be placed in an environment variable named `CLIENT_SECRET`.

Expects to find a plaintext gRPC server on `localhost:9000`, exposing the `TestApi` service (see `/proto/test.proto`)