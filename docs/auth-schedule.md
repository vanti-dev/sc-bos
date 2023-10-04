Authentication Schedule
=====================================

This document presents several authentication flows that can be used to identify an entity who wants to access Smart
Core. It enumerates the following properties

- **Who** - which entities (people or machines) will use this flow?
- **Why** - for what purpose does the subject need to access Smart Core
- **Where** - what location will be subject be in, and what network(s) will they be using
  to connect to smart core
- **How** - what does the subject need to do in order to get access?

## Apps - Active Directory Flow

- **Who?**: Landlord employees involved in monitoring or operating the building
- **Why?**
    - To access the Smart Core web apps, for:
        - Configuration & Commissioning
        - Viewing Data
        - Manual Control
- **Where** - on workstations connected to the landlord network, with direct connections
  to the app server
- **How** - by entering their domain username and password on the Keycloak login page

### Data Storage

User information is stored in an Active Directory server. The AD server will be accessible over LDAP
to allow checking user passwords, and querying basic user information including group membership.

### Technical Details

Keycloak will connect to an AD via LDAP to authenticate users.
User roles will be auto-assigned based on group memberships in Active Directory.
The Smart Core apps will send the entered username and password to Keycloak, which will return a Bearer token which
can be used for Smart Core API access.

### Management

Day-to-day management will be performed using the AD server. By assigning users to the
appropriate groups, network administrators can provide landlord employees access to the
parts of Smart Core they require.

### Reliability

This authentication flow is dependent on both Active Directory and Keycloak being operational and accessible.
Keycloak must be accessible from the client web app, so it can retrieve access tokens.
If Active Directory is not available, but Keycloak is, then users can log in with
the [Apps - Backup User Access](#apps---local-user-flow) flow below.
If Keycloak is unavailable, then users can sign in using [Apps - Local User Access](#apps---backup-user-access).
Existing sessions will continue to work until issued access tokens expire. The validity duration of tokens is
configurable in Keycloak.

## Apps - Backup User Access

- **Who?** - network administrators, or select landlord employees
- **Why?** - to access the same resources as the Active Directory flow, when the LDAP
  connection to AD is unavailable
- **Where?** - as in the AD flow
- **How?** - by entering a Smart Core-specific username and password into Keycloak the login page

### Data Storage

In addition to using user accounts via LDAP, Keycloak also maintains its own user database. Users can be created with
usernames, emails, names (and other expected profile stuff), and roles. Keycloak's user database will be stored in
the PostgreSQL database.

### Technical Details

As in "Apps - Active Directory" flow, but user is stored in Keycloak's own database rather than Active Directory.
Keycloak issues access tokens as normal.

### Management

Users must be manually provisioned using the Keyclock web interface. It is expected that only a small number of accounts
will be required.

### Reliability

Flow is available only if Keycloak is operational and accessible.

### Alternatives

Rather than having Keycloak query Active Directory for logins, it is possible to sync the directory into Keycloak's own
user database. This means that logins against directory user accounts would be possible even if Active Directory
were temporarily unavailable, so manually managing keycloak accounts would not be necessary. However, doing this
would involve storing personal data and hashed user passwords in Keycloak.

## Apps - Local User Flow

- **Who?** - network administrators, commissioners, or select landlord employees
- **Why?** - to setup or troubleshoot individual controllers, or in cases where other more robust mechanisms are
  unavailable
- **Where?** - as in the AD flow
- **How?** - by entering a Smart Core-specific username and password into the app hosted login page

### Data Storage

Users are configured via a configuration file stored on the local filesystem of the controller.
User information and claims are configured in the file and the credentials are hashed.

### Technical Details

The app invokes the OAuth Password flow to exchange user entered credentials for an access token directly with the Smart
Core controller.

### Management

Users must be manually provisioned by an administrator during the setup of the controller.

### Reliability

If the server is online, the flow is available.

### Alternatives

This is the lowest level of user authentication, Keycloak auth is preferred.
The only other alternative is to disable authentication entirely, which is not recommended.

## Tenant Token Flow

- **Who?** - Tenant's own software services (machine, not a person)
- **Why?** - To allow the tenant to monitor and control the smart systems within their demise,
  using the Smart Core API.
- **Where?** - on the tenant's own in-building network; connections go via the gateway
- **How?** - tenant exchanges their client ID and secret with the server for an access token

### Data Storage

Tenant details and secrets are stored in PostgreSQL.
Each tenant is represented in the DB and linked to one or more secrets.
The building controller it typically uniquely responsible for DB access via it's API.

### Technical Details

The tenant's client will perform an OAuth2 Client Credentials grant using their
client secret. They will receive an access token which will permit them access to the
Smart Core API via the gateway.

The API Gateway will ask the building controller to verify the token,
and will only forward the request onward if the token is valid
and the request passes all applicable security rules.

Tenant claims include zones which describe while devices the tenant is permitted to access.

### Management

Tenants are created and managed via the Ops App and/or the Smart Core Tenant API.
Tenant secrets are generated and managed via the same mechanisms.
The building FM team or ops team are responsible for sharing the tenant Client
ID and generated secrets with the tenant.

### Reliability

Accessing the Smart Core API using tenant credentials via the gateway involves the gateway talking to the building
controller and the building controller talking to the DB. If either of these are unavailable, the flow will not work.

### Alternatives

Tenant accounts can also be configured using the [Tenant Local Flow](#tenant-local-flow) below.
The gateway can be configured to connect directly to the DB at the cost of granting permission for the gateway to access
the DB directly.

## Tenant Local Flow

- **Who?** - Tenant's own software services (machine, not a person)
- **Why?** - To allow the tenant to monitor and control the smart systems within their demise,
  using the Smart Core API.
- **Where?** - on the tenant's own in-building network; connections go via the gateway
- **How?** - tenant exchanges their client ID and secret with the server for an access token

### Data Storage

Tenants and secrets are configured via a configuration file stored on the local filesystem of the controller.
Tenant information and claims are configured in the file and the credentials are hashed.

### Technical Details

The tenant authenticates with the server using the same OAuth2 Client Credentials grant as in
the [Tenant Token Flow](#tenant-token-flow).

### Management

Tenants are created during controller setup by an administrator with access to the underlying filesystem.

### Reliability

If the server is online, the flow is available.

### Alternatives

This is the lowest level of tenant authentication without any management opportunity.
