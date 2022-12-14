Installing the system for development
=====================================
This guide explains how to set up the system on your development machine. This guide DOES NOT produce a secure system
and must not be used for installation in production. 

# Dependencies
The system depends on two third-party servers:
  - Keycloak
  - PostgreSQL

These must be installed and running in order to test the system locally.

## Using Docker Compose

In the root of the repo is a `docker-compose.yml`. Running this will start Keycloak and PostgreSQL, but some manual
configuration will still be required.

# Setup

## PostgreSQL

1. Login to [pg-admin](http://localhost:8881) using the username and password defined in the docker-compose file (
   admin@example.com/postgres)
2. Create a connection to the local postgres server. This needs to be done using the IP of the docker bridge network
   rather than localhost.

   It might be `172.18.0.4`, but if that doesn't work you can find the correct IP by doing:
   `docker network inspect bsp-ew_default` and looking for `bsp-ew-db-1`

_The following steps will have been applied automatically if you used the above docker-compose file:_
   
3. Create a database called `keycloak`
4. Create a database called `smart_core`
5. `CREATE EXTENSION "uuid-ossp"` in the `smart_core` database.
6. Apply the schema from `scripts/schema.sql` to the `smart_core` database.

## Keycloak

1. [Login](http://localhost:8888) using the username and password defined in the docker-compose file (admin/admin)
2. Import realm `smart_core` from `manifests/keycloak/realm-smart-core.json` (the import option is on the 'add realm'
   page, available in the dropdown next to the current realm name (which is likely 'Master'))
