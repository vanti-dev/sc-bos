Smart Core Building Operating System
=================================

## System Architecture
The system functionality is distributed among multiple components.

### Building Controller
The Building Controller is responsible for management of the Smart Core installation. It is written in Go and installed in a
virtual machine on a server. Over gRPC, it exposes both the Smart Core API (for device control and data collection) and
private APIs for communicating with the frontend and other nodes.

It provides the following services:
  - Hosting the building user interfaces
  - Managing the enrollment of Area Controllers
    - Certificate signing
    - Smart Core CA
  - Generating & distributing configuration data to other components
  - Generating & distributing routing tables, which map device names to Smart Core nodes
  - Tenant management
    - Issuing tenant API keys
    - Assigning tenants to areas
    - Tenant data API
  - Performing whole-building automations
  - Alarm engine

### Keycloak
Smart Core requires an OpenID Connect server to provide identity & access management functionality.
We recommend [Keycloak](https://www.keycloak.org/), an open source Identity and Access Management server.
Keycloak can use an internal user database or back off to an external source such as Microsoft AD.

Keycloak hosts its own user interface (which is used as a login page), and issues OpenID Connect tokens which can be
verified by the other components, in order to provide authentication.

It is hosted in a virtual machine on a server.

### Database Server
The system uses the PostgreSQL RDBMS for persistent data storage.
Only the App Server and Keycloak directly access the database - if any other system component requires persistent data
storage, it must do so via an API exposed by the App Server.

It is hosted in a virtual machine on a server.

### Gateway
The Gateway bridges the building internal network with the untrusted networks, such as a tenant network. 
It provides a policy enforcement point to apply extra security rules to tenant API access. It also acts as an OAuth2 
server, separate from Keycloak, to issue access tokens to machines (not people).

It is hosted in a virtual machine on a server.

### Area Controllers
The Area Controller is responsible for control of devices in its local area, such as a floor. Each one runs
the same Area Controller service, but with different configuration. Configuration is manged centrally by the App Server.

The Area Controllers provide the following services:
  - Hosting a commissioning & diagnostics user interface
  - Local automations