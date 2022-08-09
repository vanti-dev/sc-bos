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
1. Create a database called `keycloak`.
2. Create a database called `smart_core`
3. `CREATE EXTENSION "uuid-ossp"` in the `smart_core` database.
4. Apply the schema from `scripts/schema.sql` to the `smart_core` database.

## Keycloak
1. Import realm `smart_core` from `manifests/keycloak/realm-smart-core.json`
