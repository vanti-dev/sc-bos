Command: `sample-tenant`
========================

This program is a very basic example of a tenant client.
The program will use a client secret to obtain an OAuth2 access token from the server, and then use the token
to perform a call that requires authorization.

Usage:

    Usage of sample-tenant:
      -addr string
            host:port of gRPC server to call (default "localhost:23557")
      -ca string
            path to Root CA certificate(s), PEM format X.509
      -client-id string
            OAuth2 client ID
      -client-secret string
            OAuth2 client secret
      -insecure
            don't verify TLS certificates
      -token-url string
            URL of OAuth2 token endpoint (default "https://localhost:8443/oauth2/token")
