Authentication/Authorisation Schedule
=====================================

This document presents several auth flows that can be used to access Smart Core in
Enterprise Wharf. It enumerates the following properties

  - **Who** - which entities (people or machines) will use this flow?
  - **Why** - for what purpose does the subject need to access Smart Core
  - **Where** - what location will be subject be in, and what network(s) will they be using
    to connect to smart core
  - **How** - what does the subject need to do in order to get access?

## Apps - Active Directory Flow
  - **Who?**: Landlord employees involved in monitoring or operating the building
  - **Why?**
    - To access the Smart Core web apps, for:
      - Configuration & Comissioning
      - Viewing Data
      - Manual Control
  - **Where** - on workstations connected to the landlord network, with direct connections
    to the app server
  - **How** - by entering their domain username and password on the login page

### Data Storage
User information is stored in an on-prem Active Directory server. The AD server will be accessible over LDAP
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
If Active Directory is not available, but Keycloak is, then users can log in with the "Apps - Backup User Access"
flow below.
If Keycloak is unavailable, then users cannot sign in to the Smart Core apps. Existing sessions will continue to work
until issued access tokens expire. The validity duration of tokens is configurable in Keycloak.

## Apps - Backup User Access
  - **Who?** - network administrators, or select landlord employees
  - **Why?** - to access the same resources as the Active Directory flow, when the LDAP
    connection to AD is unavailable
  - **Where?** - as in the AD flow
  - **How?** - by entering a Smart Core-specific username and password into the login page

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

## Tenant Token Flow
- **Who?** - Tenant's own software services (machine, not a person)
- **Why?**
    - To allow the tenant to monitor and control the smart systems within their demise,
      using the Smart Core API.
- **Where?** - on the tenant's own in-building network; connections go via the gateway
- **How?** - tenant presents the client ID and secret they were given

### Data Storage
Tenant secrets are stored in a table in PostgreSQL, which associates each secret with a particular tenant.
We will also need a database of tenants, details TBC.

### Technical Details
The tenant's client will perform an OAuth2 Client Credentials grant using their 
client secret. They will receive an access token which will permit them access to the
Smart Core API, only via the gateway.

The API Gateway will verify the tokens, and will only forward the request onward if the token is valid
and the request passes all applicable security rules.

## Area Controllers - Local Web Interface
  - **Who?** - commissioner or administrator
  - **Why?** - to directly access an Area Controller in order to commission or troubleshoot it 
  - **Where?** - on the landlord network, with direct access to the area controller
  - **How?** - TBC