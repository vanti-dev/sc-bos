BSP Enterprise Wharf - Smart Core
=================================

This repository contains the software for the Smart Core installation at Enterprise Wharf. 

## System Architecture
The system functionality is distributed among multiple components.

### App Server
The App Server is responsible for management of the Smart Core installation. It is written in Go and installed in a
virtual machine on a server. Over gRPC, it exposes both the Smart Core API (for device control and data collection) and
custom project-specific APIs.

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
  - Aggregating & reporting for the Emergency Lighting system
  - Alarm engine

### Keycloak
[Keycloak](https://www.keycloak.org/) is an open source Identity and Access Management server.
It is used at Enterprise Wharf to provide authentication and authorization services to Smart Core.
Keycloak is connected to the on-premise Microsoft Active Directory server using LDAP, which allows users of the web
apps to sign in using their domain credentials.

Keycloak hosts its own user interface (which is used as a login page), and issues OpenID Connect tokens which can be
verified by the other components, in order to provide authentication.

It is hosted in a virtual machine on a server.

### Database Server
The system uses the PostgreSQL RDBMS for persistent data storage.
Only the App Server and Keycloak directly access the database - if any other system component requires persistent data
storage, it must do so via an API exposed by the App Server.

It is hosted in a virtual machine on a server.

### Gateway
The Gateway bridges the building internal network with the untrusted tenant network. It provides a policy enforcement
point to apply extra security rules to tenant API access. It also acts as an OAuth2 server, separate from Keycloak
(which tenants can't access anyway), to issue access tokens to machines (not people).

It is hosted in a virtual machine on a server.

### Area Controllers
Each floor of the building has a [Beckhoff CX5230](https://www.beckhoff.com/en-gb/products/ipc/embedded-pcs/cx5200-intel-atom-x/cx5230.html) 
Embedded PC as an Area Controller. The Area Controller is responsible for control of devices on its floor. Each one runs
the same Area Controller service, but with different configuration. Configuration is manged centrally by the App Server.

The Embedded PCs have IO slices attached, both directly and via EtherCAT. These IO slices can only be directly accessed
via TwinCAT PLC, so they run PLC programs that can communicate with the Area Controller service over TwinCAT 3 ADS.
Lighting control via DALI happens through this mechanism.

The Area Controllers provide the following services:
  - Hosting a commissioning & diagnostics user interface
  - Local automations
  - Lighting control
  - Emergency Lighting data collection