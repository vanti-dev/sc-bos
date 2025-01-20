Smart Core BOS permissions model
================================

## Principals
*Principal*: an entity requesting access to a resources 

- User Account - a principal that is a human user. Authenticates using a username/password or OpenID Connect.
- Service Account - a principal that is an unattended machine, acting on its own behalf.
  Authenticates using a client ID and secret.
  - Touch Panel Account - a service account for a user interface that people use anonymously
  - API Account - a service account for an external Smart Core API consumer
- Node - a Smart Core BOS node in the same cohort as the node being accessed.
  Authenticates using a TLS Client Certificate signed by the Cohort Root CA.

## Resources

### Named Entities
*Named Entity*: a resource that exists within the unified Smart Core namespace.
Often just called "a device" but can actually be a few different kinds of entity that all exist in the same namespace.

- Device - an entity created by a driver to represent a building device
- Zone - an entity representing a physical area within a building.
  A zone can be directly referenced in commands, which will control the space as a whole.
- Node - a Smart Core BOS node within the cohort.
- Driver (TODO) - a service that provides devices abstraction.
  Not currently in the unified namespace, but will be in the future.
- Automation (TODO) - a service that provides automation logic.
  Not currently in the unified namespace, but will be in the future.

### Alerts
*Alert*: an event raised by Smart Core BOS to indicate a potential problem in the building.

Alerts are not necessarily associated with a specific named device. 
They have some properties, all optional, that can be used to filter them:
- Floor
- Zone
- Source - name of the entity that the alert relates
- Subsystem
- Federation - the name of the node that raised the alert

## Permissions
Permission: an action (or set of related actions) that a principal is allowed to perform on a resource.

A permission may be *scopable* or *unscopable*. A scopable permission can be restricted to a subset of
resources. An unscopable permission is inherently global.

- Trait Permissions (resource: named entity). A trait permission can be applied to all traits, or to a specific trait.
  - Read (scopable) - permitted to call read-only trait methods on named entities, such as GetXxx and PullXxx
  - Write (scopable) - permitted to call read and write trait methods on named entities, such as SetXxx and PushXxx
- Alert Permissions (resource: alert)
  - Read (scopable) - access data for the alert
  - Acknowledge (scopable) - acknowledge / unacknowledge alerts
  - Admin (unscopable) - create/update/delete custom alerts
- Service Permissions (resource: service)
  - Read (scopable) - permitted to see the state and configuration of a service
  - Lifecycle (scopable) - permitted to start, stop and restart services 
  - Configure (scopable) - permitted to alter the configuration of an existing service
    - Includes Lifecycle permission
  - Write (scopable) - permitted to create, update and delete services
    - Includes Configure permission
- Account Management Permissions (resource: principal)
  - Read (scopable) - see principal's profile metadata, role assignments, etc.
    No access to account secrets.
  - Credential (scopable) - change the principal's password or client secret.
    No access to read existing secrets, these are never exposed by the system except at creation time.
  - Write (unscopable) - create, update, delete principals and role assignments.
    No access to account secrets.
  
## Roles
A role is a collection of permissions that are grouped together and named.

Smart Core has a collection of built-in roles. Custom roles can also be created.

A role is scopable if and only if all the permissions it contains are scopable.

- Admin - has all permissions. Full access to all resources.
- Commissioner - has all Trait and Service permissions. Cannot manage principals.
- Operator
  - Trait Read, Write (all traits)
  - Service Read, Lifecycle, Configure
- Viewer
  - Trait Read (all traits)

## Role Assignment
A *Role Assignment* is the association of a role with a principal.

- A principal can have multiple roles.
- A role can be assigned to multiple principals.
- A role assignment can have a scope. This limits the role to a subset of the resources.
  The role's permissions are only effective within the specified scope. If the request accessed
  a resource outside the scope, the role won't apply to the request.
- A role assignment can only have a scope if the role is scopable. The system shall prevent the
  creation of a role assignment with a scope if the role is unscopable.
- If a role assignment has no scope, it applies to all resources.

## Scopes
A scope is a set of resources that a role assignment applies to.

### Scope kinds for named entities

To be implemented as part of MVP
- Zone Text - matches resources linked to a zone with give name, ignoring case
  - For entities, uses `metadata.location.zone` property
  - For alerts, uses the `zone` property
- Floor Text - matches all entities whose floor property is equal to the scope value, ignoring case
  - For entities, uses `metadata.location.floor` property
  - For alerts, uses the `floor` property
- Name - matches the entity whose name is equal to the scope value
  - For entities, the name of the entity is used
  - For alerts, the `source` property is used
- Name Prefix - matches all entities whose name is equal to the scope value, or starts with the scope value followed
  by a forward slash (representing a subpath)
  - e.g. `np:ns/foo` matches `ns/foo` and `ns/foo/bar` but not `ns/foobar`
- Node - the scope value is the node name. Matches all entities announced by that node.
  TODO: we need to add some device properties

To be implemented at a later date
- Group - matches all entities that are part of the group with the given ID
- Zone ID - matches all entities that are inside the zone with the given ID, based on the building model

### Scope kinds for alerts
- Zone Text - matches alerts with zone property equal to given zone, ignoring case
- Floor Text - matches alerts with floor property equal to given floor, ignoring case 
- Name - matches alerts with source property equal to given name
- Name Prefix - matches alerts with source property equal to given name, or starting with the name followed by a forward slash
- Node - matches alerts with federation property equal to given node

### Scope kinds for services
- Name
  - Matches services attached to the node with the given name
  - Matches drivers attached to the zone with the give name
- Name Prefix
  - Matches services attached to the node with the given name, or starting with the name followed by a forward slash
  - Matches drivers attached to the zone with the give name, or starting with the name followed by a forward slash
- Node
  - Matches services attached to the node with the given name
  - Matches drivers attached to a zone on the node with the given name

### Scope kinds for principals

Only supports a specific principal ID.

This is to facilitate a user being able to see / edit their own profile and change their own password.

Admin users will be given the unscoped permissions.

Note that there are no permissions to allow anyone to retrieve a principal's secrets.

## Implicit permissions
Based on the type of principal, the system will automatically grant certain permissions.

### User Accounts
Implicit *Account Management - Read* and *Account Management - Credential* scoped
to the principal's own ID.

## Deciding Access to a Resource
When a principal attempts to access a resource, it provides a credential containing authentication and authorization data.

### Access Control for Node Principals
Node principals are authenticated using a TLS Client Certificate.
At present, all Nodes are implicitly trusted with full access. Therefore, if the TLS Client Certificate validates
successfully against the Cohort Root CA, the request is allowed.

### Access Control for Other Principals
Other principals are authenticated using an OAuth2 Access Token.
We use stateless JWTs as access tokens. The token payload contains all the information needed to make an access control decision.

We extract a list of (permission, scope) pairs from the token.

Access is granted if, for any of the (permission, scope) pairs in the token:
- The permission covers the requested action
- The scope matches the requested resource