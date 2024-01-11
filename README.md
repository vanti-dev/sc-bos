Smart Core Building Operating System
=================================

For getting started with developing, see the [dev guide](docs/install/dev.md)

## System Architecture

The system functionality is distributed among multiple components.

### Building Controller

The Building Controller is responsible for management of the Smart Core installation. It is written in Go and installed
in a virtual machine on a server. Over gRPC, it exposes both the Smart Core API (for device control and data collection)
and private APIs for communicating with the frontend and other nodes.

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

## Using the Docker image

### 1. Create manifests repo

1. Create a new repository for your project, this should be called something like `<client>-<site>-service-manifests`
2. Create a folder for the target machine/gateway, e.g. `sc-gateway-1`
2. Create a `docker-compose.yml` file in the folder for the target machine/gateway
3. Add the following to the `docker-compose.yml` file:

```yaml
version: "3.7"
services:
   sc-gateway-1:
      image: ghcr.io/vanti-dev/sc-bos:v0.2024.1
      restart: always
      deploy:
         resources:
            limits:
               cpus: "0.50"
               memory: "500m"
            reservations:
               memory: "350m"
         restart_policy:
            condition: any
      volumes:
         - ./cfg:/cfg:ro
         - /var/data/sc-bos:/data:z
      ports:
         - "443:443"
         - "23557:23557"
         # the following ports are for driver webhooks - add as many as required for the number of devices (7000-7999 is supported)
         - "7000:7000"
         - "7001:7001"
      secrets:
         - steinel-password
         - xovis-password
      command: [ "--appconf=/cfg/appconf.json" ]

secrets:
   steinel-password:
      external: true
   xovis-password:
      external: true
```

4. Create a `cfg` folder in the same directory as the `docker-compose.yml` file
5. Create an `appconf.json` file in the `cfg` folder that contains the necessary Smart Core app config

### 2. Setup Auth

1. Login to GitHub as the `vanti-bot` user and create a personal access token for deployment:
    1. Go to...
    2. Click...

### 3. Deploy to target machine or gateway

1. SSH into the target machine or gateway
2. Clone the manifests repo
3. Add podman secrets
3. Run `podman-compose up -d` to start the container
