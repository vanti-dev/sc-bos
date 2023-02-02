// Package authn provides token based authentication for gRPC requests.
//
// There are two main categories of authn provided by this service:
//  1. User auth.
//     Used by applications a user interacts with to grant that application permission to act on the users behalf.
//     For example a web app that needs to execute RPCs against this server to access or modify resources owned by the user.
//  2. Machine authn.
//     Used by applications that require no user intervention. For example a cron job or the AV system of a tenant fit out.
//
// Both auth categories ultimately result in the server checking a client provided token before allowing a request to proceed.
// Authorization - deciding whether an RPC is allowed - is performed by Open Policy Agent, see docs/permissions-scheme.md
// This system contributes token validation and claims to the attributes available to OPA policies.
//
// Issuing a token for _user auth_ is either performed via a third-party authentication server or via self hosted OAuth2 Password Flow server.
// In our case the Auth Server is Keycloak which can be configured to connect to AD to authenticate the user and issue access tokens.
// The Password Flow server is written by us and checks user credentials against a json file containing ids and hashed passwords, typically used as a local account login for commissioning.
//
// _Machine to controller auth_ is similar but uses the OAuth2 Client Credentials flow.
// The client is given by an admin a Client ID and a Client Secret when their account is setup.
// Those details are then exchanged for an access token which can be used by the client to issue RPCs against the server.
// When validating the Client ID and Secret the server can use a TenantApi.VerifySecret RPC or a local JSON file containing client ids and hashed secrets.
// Tenant secrets are typically setup using the tenants system in a production environment.
package authn
